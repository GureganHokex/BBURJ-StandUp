# Comic Site (Go + Gin + GORM)

Персональный сайт комика с публичной частью и админ-панелью в стиле Django Admin.

## Стек

- Go, Gin, GORM, PostgreSQL
- HTML + Tailwind CSS (CDN) + HTMX
- Cookie-сессии, bcrypt, CSRF
- Docker Compose

## Быстрый старт (Docker, development)

```bash
cp .env.example .env
docker compose up --build
```

- Публичный сайт: http://localhost:8080
- Админка: http://localhost:8080/admin
- Логин по умолчанию задаётся в `.env` (см. `.env.example`; только для dev)

## Production и безопасность

См. [SECURITY.md](SECURITY.md).

### Yandex Cloud (бюджет ~3 200–3 500 ₽/мес)

Один сервер: приложение + PostgreSQL + Caddy (HTTPS). Без Render.

1. VM в YC: **2 vCPU, 4 GB RAM, 40 GB network-ssd**, Ubuntu 22.04, `ru-central1-a`.
2. Security Group: вход **22** (только ваш IP), **80**, **443**.
3. На VM: `sudo bash scripts/yc-bootstrap.sh`
4. `cp .env.yc.example .env`, заполните секреты и `DOMAIN`.
5. `docker compose -f docker-compose.yc.yml --env-file .env up -d --build`
6. DNS: A-запись домена → публичный IP VM.

### Docker Compose (свой сервер / VPS)

1. Скопируйте `docker-compose.prod.yml.example` → `docker-compose.prod.yml`.
2. Создайте `.env` с сильными секретами (`openssl rand -hex 32` для `SESSION_SECRET`).
3. Установите `APP_ENV=production`, `SECURE_COOKIES=true`, `DATABASE_URL` с `sslmode=require`.
4. Поставьте reverse proxy (HTTPS) и `TRUSTED_PROXIES` = IP прокси.
5. Не коммитьте `.env` в GitHub.

Миграции в production применяются при старте из `migrations/`. Вручную: `go run ./cmd/migrate`.

## Локальный запуск (без Docker)

1. Поднимите PostgreSQL и создайте БД `comic`.
2. Скопируйте `.env.example` в `.env` и поправьте `DATABASE_URL`.
3. Запустите:

```bash
go mod tidy
go run ./cmd/app
```

## Структура

- `cmd/app` — точка входа
- `internal/models`, `repository`, `services` — домен и данные
- `internal/handlers` — public, admin, api
- `web/templates` — HTML шаблоны
- `migrations/` — SQL-миграции

## API (только для авторизованной админки)

- `GET/POST /api/events`, `GET/PUT/DELETE /api/events/:id`
- `GET/POST /api/videos`, `GET/PUT/DELETE /api/videos/:id`
- `GET/POST /api/merch`, `GET/PUT/DELETE /api/merch/:id`

Мутирующие запросы требуют cookie-сессию и заголовок `X-CSRF-Token`.

## Публичные страницы

- `/` — главная
- `/events` — афиша
- `/videos` — видео (iframe)
- `/merch` — мерч

## Админка

- `/admin` — dashboard
- `/admin/settings` — тексты, соцсети, URL hero и портрета
- `/admin/account` — смена пароля администратора
- `/admin/events`, `/admin/videos`, `/admin/photos`, `/admin/merch` — CRUD через API

Цена в мерче хранится в **копейках** (например, 150000 = 1500 ₽).

**Загрузка фото:** в админке перетащите изображение в зону загрузки (Photos, Merch, Settings).
Файлы сохраняются в `uploads/` и доступны по URL `/uploads/...`.

Форматы: JPEG, PNG, WebP, GIF (до 10 МБ, настраивается `MAX_UPLOAD_MB`).
