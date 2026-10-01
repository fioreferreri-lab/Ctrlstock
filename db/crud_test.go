package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// ---------------------------------------------------------------------------
// Helpers de conexión y datos de prueba
// ---------------------------------------------------------------------------

func newTestQueries(t *testing.T) (*Queries, context.Context) {
	t.Helper()
	dbConn, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockapp?sslmode=disable")
	if err != nil {
		t.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}
	t.Cleanup(func() { dbConn.Close() })
	return New(dbConn), context.Background()
}

// crea una categoría de prueba usando la función de sqlc (no SQL a mano)
// y la borra automáticamente al terminar el test.
func createTestCategoria(t *testing.T, q *Queries, ctx context.Context, id int32) Categoria {
	t.Helper()
	cat, err := q.CreateCategoria(ctx, CreateCategoriaParams{
		IDCategoria: id,
		Nombre:      "Accesorios",
		PadreID:     sql.NullInt32{},
	})
	if err != nil {
		t.Fatalf("no se pudo preparar la categoría de prueba: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteCategoria(ctx, id) })
	return cat
}

// crea un producto de prueba y lo borra automáticamente al terminar el test.
func createTestProducto(t *testing.T, q *Queries, ctx context.Context, id, catID int32) Producto {
	t.Helper()
	p, err := q.CreateProducto(ctx, CreateProductoParams{
		IDProducto:  id,
		Nombre:      "Scrunchie Mío Mío",
		Descripcion: sql.NullString{String: "Scrunchie de seda natural", Valid: true},
		Precio:      "1500.50",
		Stock:       25,
		IDCategoria: sql.NullInt32{Int32: catID, Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo preparar el producto de prueba: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteProducto(ctx, id) })
	return p
}

// ---------------------------------------------------------------------------
// PRODUCTO
// ---------------------------------------------------------------------------

func TestCreateProducto(t *testing.T) {
	q, ctx := newTestQueries(t)
	cat := createTestCategoria(t, q, ctx, 1)

	producto, err := q.CreateProducto(ctx, CreateProductoParams{
		IDProducto:  100,
		Nombre:      "Scrunchie Mío Mío",
		Descripcion: sql.NullString{String: "Scrunchie de seda natural", Valid: true},
		Precio:      "1500.50",
		Stock:       25,
		IDCategoria: sql.NullInt32{Int32: cat.IDCategoria, Valid: true},
	})
	if err != nil {
		t.Fatalf("falló al crear producto: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteProducto(ctx, 100) })

	if producto.Nombre != "Scrunchie Mío Mío" {
		t.Errorf("esperaba nombre %q, obtuve %q", "Scrunchie Mío Mío", producto.Nombre)
	}
}

func TestGetProductoByID(t *testing.T) {
	q, ctx := newTestQueries(t)
	cat := createTestCategoria(t, q, ctx, 1)
	prod := createTestProducto(t, q, ctx, 100, cat.IDCategoria)

	producto, err := q.GetProductoByID(ctx, prod.IDProducto)
	if err != nil {
		t.Fatalf("falló al obtener producto: %v", err)
	}
	if producto.IDProducto != prod.IDProducto {
		t.Errorf("esperaba ID %d, obtuve %d", prod.IDProducto, producto.IDProducto)
	}
}

func TestGetProductoByID_Inexistente(t *testing.T) {
	q, ctx := newTestQueries(t)

	_, err := q.GetProductoByID(ctx, 99999)
	if err == nil {
		t.Errorf("se esperaba un error al buscar un ID inexistente, pero la consulta tuvo éxito")
	}
}

func TestUpdateProducto(t *testing.T) {
	q, ctx := newTestQueries(t)
	cat := createTestCategoria(t, q, ctx, 1)
	prod := createTestProducto(t, q, ctx, 100, cat.IDCategoria)

	err := q.UpdateProducto(ctx, UpdateProductoParams{
		IDProducto:  prod.IDProducto,
		Nombre:      "Scrunchie Premium",
		Descripcion: sql.NullString{String: "Edición limitada", Valid: true},
		Precio:      "1800.00",
		Stock:       30,
		IDCategoria: sql.NullInt32{Int32: cat.IDCategoria, Valid: true},
	})
	if err != nil {
		t.Fatalf("falló al actualizar producto: %v", err)
	}

	actualizado, err := q.GetProductoByID(ctx, prod.IDProducto)
	if err != nil {
		t.Fatalf("falló al releer producto actualizado: %v", err)
	}
	if actualizado.Nombre != "Scrunchie Premium" {
		t.Errorf("esperaba nombre actualizado %q, obtuve %q", "Scrunchie Premium", actualizado.Nombre)
	}
}

func TestCreateProducto_CategoriaInexistente(t *testing.T) {
	q, ctx := newTestQueries(t)

	_, err := q.CreateProducto(ctx, CreateProductoParams{
		IDProducto:  999,
		Nombre:      "Collar Roto",
		Descripcion: sql.NullString{String: "Error", Valid: true},
		Precio:      "500.00",
		Stock:       1,
		IDCategoria: sql.NullInt32{Int32: 8888, Valid: true}, // no existe
	})
	if err == nil {
		t.Errorf("se esperaba un error de clave foránea, pero el producto se creó igual")
		_ = q.DeleteProducto(ctx, 999) // limpiar si falló la aserción
	}
}

func TestDeleteProducto(t *testing.T) {
	q, ctx := newTestQueries(t)
	cat := createTestCategoria(t, q, ctx, 1)
	prod := createTestProducto(t, q, ctx, 100, cat.IDCategoria)

	if err := q.DeleteProducto(ctx, prod.IDProducto); err != nil {
		t.Fatalf("falló al eliminar producto: %v", err)
	}

	_, err := q.GetProductoByID(ctx, prod.IDProducto)
	if err == nil {
		t.Errorf("el producto debería no existir después de borrarlo")
	}
}

// ---------------------------------------------------------------------------
// MOVIMIENTO_STOCK
// ---------------------------------------------------------------------------

func TestCreateMovimiento(t *testing.T) {
	q, ctx := newTestQueries(t)
	cat := createTestCategoria(t, q, ctx, 1)
	prod := createTestProducto(t, q, ctx, 100, cat.IDCategoria)

	mov, err := q.CreateMovimiento(ctx, CreateMovimientoParams{
		IDMovimiento: 500,
		ProductoID:   prod.IDProducto,
		Tipo:         "INGRESO",
		Cantidad:     10,
		Motivo:       "Reposición inicial",
	})
	if err != nil {
		t.Fatalf("falló al crear movimiento de stock: %v", err)
	}
	t.Cleanup(func() { _ = q.DeleteMovimiento(ctx, 500) })

	if mov.Cantidad != 10 {
		t.Errorf("esperaba cantidad 10, obtuve %d", mov.Cantidad)
	}
}

func TestCreateMovimiento_ProductoInexistente(t *testing.T) {
	q, ctx := newTestQueries(t)

	_, err := q.CreateMovimiento(ctx, CreateMovimientoParams{
		IDMovimiento: 501,
		ProductoID:   99999, // no existe
		Tipo:         "EGRESO",
		Cantidad:     5,
		Motivo:       "Venta fantasma",
	})
	if err == nil {
		t.Errorf("se esperaba un error de clave foránea al asociar un movimiento a un producto inexistente")
		_ = q.DeleteMovimiento(ctx, 501)
	}
}
