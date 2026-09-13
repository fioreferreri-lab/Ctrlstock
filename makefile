APP_NAME := Ctrlstock
DB_URL := postgres://user:password@localhost:5432/stockapp?sslmode=disable
SCHEMA_FILE := db/scheme/schema.sql
CONTAINER_NAME := mi_base_postgres

.PHONY: all generate build test clean

all: test

generate:
	@echo "-> Generando código con sqlc..."
	@sqlc generate

build: generate
	@echo "-> Compilando..."
	@go build ./...

test: generate
	@echo "-> Limpiando entorno (contenedores y volúmenes)..."
	@docker compose down -v
	@echo "-> Levantando PostgreSQL..."
	@docker compose up -d
	@echo "-> Esperando a que la base esté lista..."
	@sleep 5
	@echo "-> Aplicando schema.sql..."
	@docker exec -i $(CONTAINER_NAME) psql -U user -d stockapp < $(SCHEMA_FILE)
	@echo "-> Corriendo tests..."
	@go test -v ./...
	@echo "-> Limpiando entorno (contenedores y volúmenes)..."
	@docker compose down -v

clean:
	@docker compose down -v
	@rm -rf tmp