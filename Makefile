.PHONY: build migrate-up migrate-down install-migrate

APP_NAME=myapp
N ?= 1

build:
	go build -o $(APP_NAME)

install-migrate:
	go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2

migrate-up:
	migrate -source file://migrations -database pgx5://admin:root@localhost:5432/project_a?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database pgx5://admin:root@localhost:5432/project_a?sslmode=disable down $(N)