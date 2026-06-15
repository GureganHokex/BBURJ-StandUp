FROM golang:1.23-alpine AS builder

# Alpine CDN is often slow/blocked from RU VPS — use Yandex mirror.
RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|https://mirror.yandex.ru/mirrors/alpine|g' /etc/apk/repositories

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,https://goproxy.io,direct

RUN apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build admin Tailwind CSS (no CDN at runtime).
RUN ARCH="$(uname -m)" \
	&& case "$ARCH" in \
		x86_64) TW_ARCH="x64" ;; \
		aarch64) TW_ARCH="arm64" ;; \
		*) echo "unsupported arch: $ARCH" && exit 1 ;; \
	esac \
	&& curl -fsSL "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-${TW_ARCH}" -o /tmp/tailwindcss \
	&& chmod +x /tmp/tailwindcss \
	&& /tmp/tailwindcss -i web/static/css/admin-input.css -o web/static/css/tailwind.css --minify

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/comic ./cmd/app

FROM alpine:3.20

RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|https://mirror.yandex.ru/mirrors/alpine|g' /etc/apk/repositories \
	&& apk add --no-cache ca-certificates tzdata \
	&& adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=builder /bin/comic /app/comic
COPY migrations /app/migrations
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app/uploads /app/migrations

USER appuser
EXPOSE 8080
CMD ["/app/comic"]
