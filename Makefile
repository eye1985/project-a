.PHONY: build migrate-up migrate-down install-migrate docker-prod docker-push-prod start-prod

APP_NAME=myapp-linux-amd-64
N ?= 1

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)

docker-prod:
	@echo "Bulding bin..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)
	@echo "Building docker img..."
	docker build . --file server.Dockerfile -t eye1985/project-a:prod --no-cache --platform linux/amd64

docker-push-prod:
	docker push eye1985/project-a:prod

start-prod:
	docker compose up --pull always --force-recreate -d

install-migrate:
	go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2

migrate-up:
	migrate -source file://migrations -database $(PGX5_URL) up

migrate-down:
	migrate -source file://migrations -database $(PGX5_URL) down $(N)