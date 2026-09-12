package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // Driver de PostgreSQL
)

func TestProductoQueries_CRUD(t *testing.T) {
	// 1. Establecer conexión a la base de datos de prueba
	dbConn, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockapp?sslmode=disable")
	if err != nil {
		t.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer dbConn.Close()

	// 2. Instanciar el objeto Queries generado por sqlc
	queries := New(dbConn)
	ctx := context.Background()

	// Variable para almacenar el ID del producto creado
	var createdProductoID int32

	// --- 1. CREATE ---
	t.Run("CreateProducto", func(t *testing.T) {
		params := CreateProductoParams{
			Nombre: "Monitor 24 pulgadas",
			// Para campos TEXT/VARCHAR que admiten nulos, Go usa sql.NullString
			Descripcion: sql.NullString{String: "Monitor Full HD", Valid: true},
			// Los campos DECIMAL muchas veces se mapean como strings para no perder precisión
			Precio: "250000.50",
			Stock:  15,
		}

		producto, err := queries.CreateProducto(ctx, params)
		if err != nil {
			t.Fatalf("Error al crear el producto: %v", err)
		}

		if producto.Nombre != params.Nombre {
			t.Errorf("Se esperaba el nombre %s, pero se obtuvo %s", params.Nombre, producto.Nombre)
		}

		createdProductoID = producto.ID // Guardamos el ID para usarlo luego
	})

	// --- 2. READ (GET) ---
	t.Run("GetProductoByID", func(t *testing.T) {
		producto, err := queries.GetProductoByID(ctx, createdProductoID)
		if err != nil {
			t.Fatalf("Error al obtener el producto: %v", err)
		}
		if producto.ID != createdProductoID {
			t.Errorf("Se esperaba el ID %d, pero se obtuvo %d", createdProductoID, producto.ID)
		}
	})

	// --- 3. UPDATE ---
	t.Run("UpdateProducto", func(t *testing.T) {
		params := UpdateProductoParams{
			ID:          createdProductoID,
			Nombre:      "Monitor 27 pulgadas",
			Descripcion: sql.NullString{String: "Monitor 4K", Valid: true},
			Precio:      "350000.00",
			Stock:       10,
		}

		err := queries.UpdateProducto(ctx, params)
		if err != nil {
			t.Fatalf("Error al actualizar el producto: %v", err)
		}
	})

	// --- 4. LIST ---
	t.Run("ListProductos", func(t *testing.T) {
		productos, err := queries.ListProductos(ctx)
		if err != nil {
			t.Fatalf("Error al listar productos: %v", err)
		}

		if len(productos) == 0 {
			t.Errorf("Se esperaba al menos 1 producto en la lista")
		}
	})

	// --- 5. DELETE ---
	t.Run("DeleteProducto", func(t *testing.T) {
		err := queries.DeleteProducto(ctx, createdProductoID)
		if err != nil {
			t.Fatalf("Error al eliminar el producto: %v", err)
		}

		// Verificamos que ya no exista
		_, err = queries.GetProductoByID(ctx, createdProductoID)
		if err == nil {
			t.Errorf("Se esperaba un error al buscar un producto eliminado, pero no se obtuvo ninguno")
		}
	})
}
