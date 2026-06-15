# Security

## Secrets and GitHub

- Never commit `.env`, keys, or database dumps.
- Use GitHub **Secrets** for CI/CD (`SESSION_SECRET`, `DATABASE_URL`, etc.).
- Rotate `SESSION_SECRET` if it may have leaked (all users will need to log in again).

## Production environment

| Variable | Requirement |
|----------|-------------|
| `APP_ENV` | `production` |
| `SESSION_SECRET` | ≥ 32 random characters, not the example value |
| `ADMIN_PASSWORD` | ≥ 12 characters, not `admin123` (only used when DB has no admin user) |
| `SECURE_COOKIES` | `true` (HTTPS) |
| `DATABASE_URL` | `sslmode=require` or stricter |
| `TRUSTED_PROXIES` | IP/CIDR of reverse proxy (nginx, Caddy) |

The app **refuses to start** if production config is unsafe.

## Deployment

1. Terminate TLS at nginx/Caddy; proxy to `app:8080`.
2. Do not expose PostgreSQL port publicly.
3. Persist `uploads/` volume and back it up with the database.
4. Run migrations: automatic on start in production, or `go run ./cmd/migrate`.

## Implemented protections

- bcrypt passwords
- HttpOnly session cookies, SameSite=Lax
- Sessions stored in PostgreSQL (survive restarts)
- CSRF on admin login and API mutations
- Open redirect protection on post-login `next`
- Rate limits: login (10 / 15 min per IP), uploads (60 / hour per IP)
- Upload: size limit, MIME sniffing, random filenames
- Security headers (strict CSP, X-Frame-Options, HSTS in production)
- External URLs restricted to `http`/`https` schemes
- Admin password change at `/admin/account`
- API pagination cap (100)
- Generic 500 errors in production (no stack traces to clients)

## Reporting issues

If you find a vulnerability, contact the maintainer privately before public disclosure.
