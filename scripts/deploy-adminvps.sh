#!/usr/bin/env bash
# Деплой на AdminVPS (или любой Ubuntu VPS) с локальной машины.
#
# Usage:
#   ./scripts/deploy-adminvps.sh <SERVER_IP> [DOMAIN]
#
# Пример:
#   ./scripts/deploy-adminvps.sh 185.12.34.56 bigburjprodj.com
#
# Требования: SSH-доступ root@IP (ключ или пароль), rsync, ssh.

set -euo pipefail

SERVER_IP="${1:-}"
DOMAIN="${2:-bigburjprodj.com}"
SSH_USER="${SSH_USER:-root}"
REMOTE_DIR="${REMOTE_DIR:-/opt/bburj-standup}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if [[ -z "${SERVER_IP}" ]]; then
	echo "Usage: $0 <SERVER_IP> [DOMAIN]"
	exit 1
fi

SSH_TARGET="${SSH_USER}@${SERVER_IP}"

echo "==> Проверка SSH: ${SSH_TARGET}"
ssh -o ConnectTimeout=15 -o StrictHostKeyChecking=accept-new "${SSH_TARGET}" "echo OK"

echo "==> Копирование проекта в ${REMOTE_DIR}"
ssh "${SSH_TARGET}" "mkdir -p ${REMOTE_DIR}"
rsync -az --delete \
	--exclude '.git' \
	--exclude '.env' \
	--exclude 'uploads' \
	--exclude 'bin' \
	--exclude '.DS_Store' \
	"${REPO_ROOT}/" "${SSH_TARGET}:${REMOTE_DIR}/"

echo "==> Bootstrap (Docker, firewall)"
ssh "${SSH_TARGET}" "bash ${REMOTE_DIR}/scripts/yc-bootstrap.sh"

echo "==> Освобождение портов 80/443 (ISPmanager/nginx на Promo)"
ssh "${SSH_TARGET}" 'for s in nginx apache2 ihttpd; do systemctl stop "$s" 2>/dev/null || true; systemctl disable "$s" 2>/dev/null || true; done'

echo "==> Генерация .env"
POSTGRES_PASSWORD="$(openssl rand -hex 16)"
SESSION_SECRET="$(openssl rand -hex 32)"
ADMIN_PASSWORD="$(openssl rand -base64 18 | tr -d '/+=' | head -c 16)"

ssh "${SSH_TARGET}" "cat > ${REMOTE_DIR}/.env" <<EOF
DOMAIN=${DOMAIN}
POSTGRES_USER=comic
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
POSTGRES_DB=comic
DATABASE_URL=postgres://comic:${POSTGRES_PASSWORD}@db:5432/comic?sslmode=prefer
SESSION_SECRET=${SESSION_SECRET}
ADMIN_USERNAME=admin
ADMIN_PASSWORD=${ADMIN_PASSWORD}
TRUSTED_PROXIES=172.16.0.0/12
EOF

echo "==> Сборка и запуск контейнеров"
ssh "${SSH_TARGET}" "cd ${REMOTE_DIR} && docker compose -f docker-compose.yc.yml --env-file .env up -d --build"

echo ""
echo "=============================================="
echo "Деплой завершён: https://${DOMAIN}"
echo "Админка: https://${DOMAIN}/admin"
echo "Логин: admin"
echo "Пароль админки (сохраните): ${ADMIN_PASSWORD}"
echo ""
echo "DNS: A-запись ${DOMAIN} -> ${SERVER_IP}"
echo "Проверка: curl -sS https://${DOMAIN}/health"
echo "=============================================="
