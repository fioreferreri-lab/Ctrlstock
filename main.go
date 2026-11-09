package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	db "ctrlstock/db"
	"ctrlstock/logic"
	"ctrlstock/views"
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

	queries := db.New(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", indexHandler(queries))
	mux.HandleFunc("POST /productos", crearProductoHandler(queries))

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func indexHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productos, err := queries.ListProductos(r.Context())
		if err != nil {
			log.Printf("ListProductos: %v", err)
			http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
			return
		}
		if err := views.IndexPage("CtrlStock - Productos", productos).Render(r.Context(), w); err != nil {
			log.Printf("Render IndexPage: %v", err)
			http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
		}
	}
}

func crearProductoHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}

		nombre := strings.TrimSpace(r.FormValue("nombre"))

		precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
		if err != nil {
			http.Error(w, "El precio debe ser un número", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			http.Error(w, "El stock debe ser un número entero", http.StatusBadRequest)
			return
		}

		idCategoria, err := nullInt32(r.FormValue("id_categoria"))
		if err != nil {
			http.Error(w, "El ID de categoría debe ser un número entero", http.StatusBadRequest)
			return
		}

		if err := logic.ValidateProducto(nombre, precio, int32(stock)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = queries.CreateProducto(r.Context(), db.CreateProductoParams{
			Nombre:      nombre,
			Descripcion: nullString(r.FormValue("descripcion")),
			Precio:      fmt.Sprintf("%.2f", precio),
			Stock:       int32(stock),
			IDCategoria: idCategoria,
		})
		if err != nil {
			log.Printf("CreateProducto: %v", err)
			http.Error(w, "No se pudo guardar el producto", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func nullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{String: s, Valid: s != ""}
}

func nullInt32(s string) (sql.NullInt32, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullInt32{}, nil
	}
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return sql.NullInt32{}, err
	}
	return sql.NullInt32{Int32: int32(n), Valid: true}, nil
}
