package main

import (
	"database/sql"
	"fmt"
	"log"

	db "sistema_stock/db" // Apuntamos a la carpeta exacta donde están los archivos de sqlc

	_ "github.com/lib/pq" // Driver de conexión
)

func main() {
	// 1. Abrimos la conexión
	dsn := "postgres://user:password@localhost:5432/stockapp?sslmode=disable"
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error configurando la conexión: ", err)
	}
	defer conn.Close()

	// Verificamos que Docker y PostgreSQL estén respondiendo
	if err := conn.Ping(); err != nil {
		log.Fatal("La base no responde: ", err)
	}
	fmt.Println("¡Conexión exitosa a PostgreSQL! 🎉")

	// 2. Inicializamos sqlc pasándole la conexión para dejar todo listo para el próximo práctico
	// Usamos "_" para instanciar el motor sin que Go tire error de "variable no usada"
	_ = db.New(conn)

	fmt.Println("✅ Paquete de base de datos (db) inicializado y listo para usar.")
}
