package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	db "ctrlstock/db"
	"ctrlstock/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s desde %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func main() {

	/*====================================
	CONEXION BASE DE DATOS
	=====================================*/
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Falta la variable de entorno DATABASE_URL")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión: %v", err)
	}
	defer conn.Close()

	var pingErr error
	for i := 0; i < 30; i++ {
		if pingErr = conn.Ping(); pingErr == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if pingErr != nil {
		log.Fatalf("No se pudo conectar a la base de datos tras varios intentos: %v", pingErr)
	}

	server := &handlers.Server{
		Queries: db.New(conn),
	}

	/*====================================
	END-POINTS
	=====================================*/
	mux := http.NewServeMux()
	mux.Handle("/categoria", loggingMiddleware(http.HandlerFunc(server.CategoriasHandler)))
	mux.Handle("/categoria/", loggingMiddleware(http.HandlerFunc(server.CategoriasHandler)))
	mux.Handle("/producto", loggingMiddleware(http.HandlerFunc(server.ProductosHandler)))
	mux.Handle("/producto/", loggingMiddleware(http.HandlerFunc(server.ProductosHandler)))
	mux.Handle("/movimiento_stock", loggingMiddleware(http.HandlerFunc(server.MovimientosHandler)))
	mux.Handle("/movimiento_stock/", loggingMiddleware(http.HandlerFunc(server.MovimientosHandler)))

	port := ":8080"
	log.Printf("Servidor escuchando en http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Error al iniciar el servidor: %s", err)
	}
}
