FROM golang:1.23-alpine AS builder

# Alpine CDN is often slow/blocked from RU VPS — use Yandex mirror.
RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|https://mirror.yandex.ru/mirrors/alpine|g' /etc/apk/repositories

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,https://goproxy.io,direct

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/comic ./cmd/app

FROM alpine:3.20

RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|https://mirror.yandex.ru/mirrors/alpine|g' /etc/apk/repositories \
	&& apk add --no-cache ca-certificates tzdata su-exec \
	&& adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=builder /bin/comic /app/comic
COPY migrations /app/migrations
COPY deploy/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh \
	&& mkdir -p /app/uploads \
	&& chown -R appuser:appuser /app/uploads /app/migrations

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/comic"]
