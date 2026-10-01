package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Driver de PostgreSQL necesario para database/sql

	"ctrlstock/db"       // Paquete generado por sqlc
	"ctrlstock/handlers" // Paquete de manejadores HTTP
)

func main() {
	// 1. Conexión a la base de datos PostgreSQL
	dsn := "postgres://user:password@localhost:5432/stockapp?sslmode=disable"
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error configurando la conexión: ", err)
	}
	defer conn.Close()

	// Verificamos que PostgreSQL/Docker responda
	if err := conn.Ping(); err != nil {
		log.Fatal("La base de datos no responde: ", err)
	}
	fmt.Println("¡Conexión exitosa a PostgreSQL!")

	// 2. Inicialización de sqlc
	queries := db.New(conn)
	_ = queries // Evita el error 'declared and not used' mientras no se use en los handlers
	fmt.Println("Paquete de base de datos (db) inicializado y listo para usar.")

	// 3. Configuración del ruteador HTTP
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/categories/", handlers.CategoryHandler)

	// 4. Inicio del servidor HTTP
	fmt.Println("Servidor de Categorías corriendo en http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
