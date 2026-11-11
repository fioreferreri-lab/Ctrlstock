package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"ctrlstock/db"
	"ctrlstock/handlers" // Importamos carpeta de handlers
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Falta la variable de entorno DATABASE_URL")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión: %v", err)
	}
	defer conn.Close()

	// Wait-for-it para Docker
	var pingErr error
	for i := 0; i < 30; i++ {
		if pingErr = conn.Ping(); pingErr == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if pingErr != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", pingErr)
	}

	// 1. Inicializamos las queries generadas por sqlc
	queries := db.New(conn)

	// 2. Instanciamos el Server de tu paquete handlers
	app := &handlers.Server{
		Queries: queries,
		DB:      conn,
	}

	mux := http.NewServeMux()

	// 3. Enrutamiento: delegamos todo a la carpeta handlers
	mux.HandleFunc("/producto", app.ProductosHandler)
	mux.HandleFunc("/producto/", app.ProductosHandler)

	mux.HandleFunc("/categoria", app.CategoriasHandler)
	mux.HandleFunc("/categoria/", app.CategoriasHandler)

	mux.HandleFunc("/movimiento_stock", app.MovimientosHandler)
	mux.HandleFunc("/movimiento_stock/", app.MovimientosHandler)

	// Redirección de la raíz al inventario de productos (Opcional, pero útil)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/producto", http.StatusSeeOther)
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
