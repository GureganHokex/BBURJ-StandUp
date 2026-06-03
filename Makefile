.PHONY: run build tidy docker-up docker-down migrate-up

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
