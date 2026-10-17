APP_NAME       := Ctrlstock
SCHEMA_FILE    := db/scheme/schema.sql
CONTAINER_NAME := mi_base_postgres
DB_USER        := user
DB_NAME        := stockapp
APP_URL        := http://localhost:8080

# Cadena de conexión que la app lee al arrancar (debe coincidir con docker-compose.yaml)
DATABASE_URL   := postgres://user:password@localhost:5432/stockapp?sslmode=disable
export DATABASE_URL

.PHONY: all run generate build test db-up db-wait schema db-reset clean

all: build

# Levanta la base, crea las tablas, genera el código y arranca la app
run: db-up db-wait schema generate
	@echo "-> Iniciando $(APP_NAME)..."
	@go run .

# Genera el código de las queries con sqlc
generate:
	@echo "-> Generando código con sqlc..."
	@sqlc generate

# Levanta el contenedor de PostgreSQL
db-up:
	@echo "-> Levantando PostgreSQL..."
	@docker compose up -d

# Espera a que PostgreSQL acepte conexiones
db-wait:
	@echo "-> Esperando a que la base esté lista..."
	@until docker exec $(CONTAINER_NAME) pg_isready -h localhost -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do sleep 1; done

# Crea las tablas a partir de schema.sql (usar CREATE TABLE IF NOT EXISTS)
schema:
	@echo "-> Aplicando schema.sql..."
	@docker exec -i $(CONTAINER_NAME) psql -v ON_ERROR_STOP=1 -U $(DB_USER) -d $(DB_NAME) < $(SCHEMA_FILE)

# Borra la base y todos sus datos
db-reset:
	@echo "-> Borrando base de datos..."
	@docker compose down -v

# Construye el binario de la aplicación
build: generate
	@echo "-> Compilando..."
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Tareas previas -> pruebas -> tareas posteriores
# 1) sqlc, compilación, borrar y levantar contenedores, esperar la base, levantar la app
# 2) tests de Go y prueba de cada endpoint
# 3) apagar la app y borrar contenedores y volúmenes (aunque algo falle)
test: build db-reset db-up db-wait schema
	@echo "-> Levantando la aplicación..."
	@./tmp/$(APP_NAME) > tmp/app.log 2>&1 & echo $$! > tmp/app.pid
	@echo "-> Esperando a que la aplicación responda..."
	@for i in $$(seq 1 30); do curl -s -o /dev/null $(APP_URL)/categoria && break; sleep 1; done; \
	if ! curl -s -o /dev/null $(APP_URL)/categoria; then \
		echo "!! La aplicación no arrancó. Contenido de tmp/app.log:"; \
		cat tmp/app.log; \
		kill `cat tmp/app.pid` 2>/dev/null; rm -f tmp/app.pid; \
		$(MAKE) --no-print-directory db-reset; \
		exit 1; \
	fi
	@status=0; \
	echo "-> Corriendo tests de Go..."; \
	go test -count=1 -v ./... || status=1; \
	$(MAKE) --no-print-directory schema > /dev/null; \
	echo "-> Probando endpoints..."; \
	bash ./requests.bash || status=1; \
	echo "-> Apagando la aplicación..."; \
	kill `cat tmp/app.pid` 2>/dev/null; rm -f tmp/app.pid; \
	$(MAKE) --no-print-directory db-reset; \
	exit $$status

# Limpia la base y los artefactos de construcción
clean: db-reset
	@rm -rf tmp