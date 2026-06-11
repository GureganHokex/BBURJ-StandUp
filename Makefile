.PHONY: run build tidy docker-up docker-down migrate-up yc-up yc-down

yc-up:
	docker compose -f docker-compose.yc.yml --env-file .env up -d --build

yc-down:
	docker compose -f docker-compose.yc.yml --env-file .env down

run:
	go run ./cmd/app

build:
	go build -o bin/comic ./cmd/app

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	@echo "Apply migrations/001_init.up.sql with golang-migrate or psql"
