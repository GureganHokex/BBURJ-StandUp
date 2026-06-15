# Большой Буржинский — сайт комика

Персональный сайт стендап-комика: публичная витрина и админ-панель для управления контентом.

---

## Возможности

**Публичная часть**
- Главная с афишей, блоком «Обо мне», видео, мерчем и фото
- Страницы афиши, видео, фото и мерча
- Модальное окно с подробностями события и ссылкой на покупку билетов
- Lightbox для галереи фото
- Адаптивная вёрстка под мобильные устройства

**Админ-панель** (`/admin`)
- Настройки сайта: тексты, соцсети, hero-изображение, портрет
- CRUD для событий, видео, фото и мерча
- Импорт афиши из билетных агрегаторов (TicketsCloud, Timepad)
- Предпросмотр постера и карточки события по ссылке на билеты
- Загрузка изображений drag-and-drop
- Смена пароля администратора

---

## Стек технологий

| Слой | Технологии |
|------|------------|
| Backend | Go 1.23, Gin, GORM |
| База данных | PostgreSQL 16 |
| Шаблоны | `html/template`, встроенные через `embed.FS` |
| Публичный UI | HTML, CSS, vanilla JavaScript |
| Админка | HTML, Tailwind CSS 3 (сборка при Docker build), vanilla JavaScript |
| Безопасность | Cookie-сессии, bcrypt, CSRF, rate limiting |
| Инфраструктура | Docker, Docker Compose, Caddy (HTTPS) |

---

## Структура репозитория

```
cmd/app/              — точка входа приложения
internal/
  config/             — конфигурация и валидация окружения
  database/           — подключение к БД, миграции
  handlers/           — HTTP: public, admin, api
  middleware/         — auth, CSRF, security headers, rate limit
  models/             — модели GORM
  repository/         — слой доступа к данным
  services/           — бизнес-логика
  session/            — хранилище сессий
  storage/            — загрузка файлов
  tickets/            — клиенты билетных агрегаторов
migrations/           — SQL-миграции (production)
web/
  templates/          — HTML-шаблоны (public + admin)
  static/             — CSS, JS, изображения
deploy/               — конфигурация Caddy
scripts/              — вспомогательные скрипты
```

---

## Публичные страницы

| URL | Описание |
|-----|----------|
| `/` | Главная |
| `/events` | Афиша выступлений |
| `/videos` | Видео (YouTube, Rutube) |
| `/photos` | Фото-галерея |
| `/merch` | Мерч |

---

## Админка

| URL | Описание |
|-----|----------|
| `/admin` | Dashboard |
| `/admin/settings` | Тексты, соцсети, изображения |
| `/admin/events` | События и афиша |
| `/admin/videos` | Видео |
| `/admin/photos` | Фото |
| `/admin/merch` | Мерч |
| `/admin/account` | Смена пароля |

**Мерч:** цена хранится в **копейках** (150000 = 1500 ₽).

**Загрузка файлов:** JPEG, PNG, WebP, GIF — до 10 МБ (`MAX_UPLOAD_MB`). Файлы доступны по `/uploads/...`.

---

## API

JSON API доступен только авторизованному администратору (cookie-сессия).

| Ресурс | Эндпоинты |
|--------|-----------|
| События | `GET/POST /api/events`, `GET/PUT/DELETE /api/events/:id`, `GET /api/events/preview-ticket` |
| Видео | `GET/POST /api/videos`, `GET/PUT/DELETE /api/videos/:id` |
| Фото | `GET/POST /api/photos`, `GET/PUT/DELETE /api/photos/:id` |
| Мерч | `GET/POST /api/merch`, `GET/PUT/DELETE /api/merch/:id` |
| Настройки | `GET/PUT /api/settings` |
| Загрузка | `POST /api/upload` |
| Афиша (импорт) | `GET /api/ticket-catalog/providers`, `GET /api/ticket-catalog/:source/events` |
| Аккаунт | `PUT /api/account/password` |

Мутирующие запросы (`POST`, `PUT`, `DELETE`) требуют заголовок `X-CSRF-Token` (токен в `<meta name="csrf-token">` админки).

Health-check: `GET /health` → `ok`.

---

## Безопасность

Требования к production-окружению, список встроенных защит и правила работы с секретами — в [SECURITY.md](SECURITY.md).
