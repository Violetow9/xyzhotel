BINARY_NAME=hotel-app

DB_DSN="root:root@tcp(127.0.0.1:3306)/xyzhotel?parseTime=true"

.PHONY: all build clean run run-cli test docker-up docker-down sqlc help

all: build

build:
	@echo "Compiling application..."
	go build -o $(BINARY_NAME) main.go
	@echo "Compiled ./$(BINARY_NAME)"

clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
	@echo "Cleaned."

run:
	@echo "Running app with http server..."
	DB_DSN=$(DB_DSN) go run main.go

run-cli:
	@echo "Running app with http server and CLI..."
	DB_DSN=$(DB_DSN) go run main.go cli

start: build
	@echo "Starting application..."
	DB_DSN=$(DB_DSN) ./$(BINARY_NAME)

docker-db:
	docker compose up -d mysql

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-reset:
	@echo "Resetting database..."
	docker compose down -v
	docker compose up -d mysql
	@echo "Database reset complete."

sqlc:
	sqlc generate

test:
	go test ./...

help:
	@echo "Commandes disponibles :"
	@echo "  make build       : Compile l'application"
	@echo "  make clean       : Supprime le binaire"
	@echo "  make run         : Lance le serveur HTTP (go run)"
	@echo "  make run-cli     : Lance le CLI interactif (go run)"
	@echo "  make docker-db   : Lance MySQL via Docker"
	@echo "  make docker-reset: Réinitialise la DB (supprime les données)"
	@echo "  make sqlc        : Génère le code Go depuis le SQL"