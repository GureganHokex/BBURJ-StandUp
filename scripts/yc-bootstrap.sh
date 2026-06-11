#!/usr/bin/env bash
# Первичная настройка Ubuntu 22.04/24.04 (VPS / Yandex Cloud). Запуск: sudo bash scripts/yc-bootstrap.sh
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
	echo "Запустите от root: sudo bash $0"
	exit 1
fi

export DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get install -y ca-certificates curl git ufw

if ! command -v docker >/dev/null 2>&1; then
	install -m 0755 -d /etc/apt/keyrings
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
	chmod a+r /etc/apt/keyrings/docker.asc
	echo \
		"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		$(. /etc/os-release && echo "${VERSION_CODENAME}") stable" \
		>/etc/apt/sources.list.d/docker.list
	apt-get update
	apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
fi

systemctl enable --now docker

ufw default deny incoming
ufw default allow outgoing
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

echo ""
echo "Docker установлен. Дальше:"
echo "  1. git clone <repo> && cd <repo>"
echo "  2. cp .env.yc.example .env && nano .env"
echo "  3. docker compose -f docker-compose.yc.yml --env-file .env up -d --build"
echo "  4. DNS: A-запись DOMAIN -> публичный IP этой VM"
