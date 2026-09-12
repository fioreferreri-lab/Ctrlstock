package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestSistema_CRUD_ValidAndInvalid(t *testing.T) {
	dbConn, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockapp?sslmode=disable")
	if err != nil {
		t.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer dbConn.Close()

	queries := New(dbConn)
	ctx := context.Background()

	// IDs fijos para las pruebas (como tus tablas no usan serial, los pasamos explícitamente)
	var catID int32 = 1
	var prodID int32 = 100
	var movID int32 = 500

	// limpieza previa correcta (esperan solo 1 valor: antes saltaba error aca)
	_ = queries.DeleteMovimiento(ctx, movID)
	_ = queries.DeleteProducto(ctx, prodID)
	// (Nota: Si tuvieras DeleteCategoria lo limpiarías acá también)

	// =========================================================================
	// 1. CATEGORÍA (Requisito previo para el producto)
	// =========================================================================
	// Nota: Como no pasaste el query de categoría en el snippet, asumimos inserción directa o la creamos por SQL plano si fuera necesario,
	// pero para este test insertamos una categoría base directamente vía SQL para cumplir la FK del producto.
	_, err = dbConn.ExecContext(ctx, "INSERT INTO categoria (id_categoria, nombre, padre_id) VALUES ($1, $2, NULL) ON CONFLICT (id_categoria) DO NOTHING", catID, "Accesorios")
	if err != nil {
		t.Fatalf("Error al preparar la categoría de prueba: %v", err)
	}

	// =========================================================================
	// 2. PRODUCTO - CASOS VÁLIDOS
	// =========================================================================
	t.Run("Validar Creación de Producto", func(t *testing.T) {
		params := CreateProductoParams{
			IDProducto:  prodID,
			Nombre:      "Scrunchie Mío Mío",
			Descripcion: sql.NullString{String: "Scrunchie de seda natural", Valid: true},
			Precio:      "1500.50",
			Stock:       25,
			IDCategoria: sql.NullInt32{Int32: catID, Valid: true},
		}

		producto, err := queries.CreateProducto(ctx, params)
		if err != nil {
			t.Fatalf("Caso válido falló al crear producto: %v", err)
		}
		if producto.Nombre != params.Nombre {
			t.Errorf("Esperaba nombre %s, obtuve %s", params.Nombre, producto.Nombre)
		}
	})

	t.Run("Validar Obtener Producto por ID", func(t *testing.T) {
		producto, err := queries.GetProductoByID(ctx, prodID)
		if err != nil {
			t.Fatalf("Caso válido falló al obtener producto: %v", err)
		}
		if producto.IDProducto != prodID {
			t.Errorf("Esperaba ID %d, obtuve %d", prodID, producto.IDProducto)
		}
	})

	t.Run("Validar Actualizar Producto", func(t *testing.T) {
		params := UpdateProductoParams{
			IDProducto:  prodID,
			Nombre:      "Scrunchie Premium",
			Descripcion: sql.NullString{String: "Edición limitada", Valid: true},
			Precio:      "1800.00",
			Stock:       30,
			IDCategoria: sql.NullInt32{Int32: catID, Valid: true},
		}

		err := queries.UpdateProducto(ctx, params)
		if err != nil {
			t.Fatalf("Caso válido falló al actualizar producto: %v", err)
		}
	})

	// =========================================================================
	// 3. PRODUCTO - CASOS INVÁLIDOS
	// =========================================================================
	t.Run("Caso Inválido - Obtener Producto Inexistente", func(t *testing.T) {
		var idInvalido int32 = 99999
		_, err := queries.GetProductoByID(ctx, idInvalido)
		if err == nil {
			t.Errorf("Se esperaba un error al buscar un ID inexistente, pero la consulta tuvo éxito")
		}
	})

	t.Run("Caso Inválido - Crear Producto con Categoría Inexistente (Falla de FK)", func(t *testing.T) {
		paramsInvalido := CreateProductoParams{
			IDProducto:  999,
			Nombre:      "Collar Roto",
			Descripcion: sql.NullString{String: "Error", Valid: true},
			Precio:      "500.00",
			Stock:       1,
			IDCategoria: sql.NullInt32{Int32: 8888, Valid: true}, // Categoría que no existe
		}

		_, err := queries.CreateProducto(ctx, paramsInvalido)
		if err == nil {
			t.Errorf("Se esperaba un error de llave foránea (Foreign Key), pero el producto se creó igual")
		}
	})

	// =========================================================================
	// 4. MOVIMIENTO DE STOCK - CASOS VÁLIDOS E INVÁLIDOS
	// =========================================================================
	t.Run("Validar Creación de Movimiento de Stock", func(t *testing.T) {
		paramsMov := CreateMovimientoParams{
			IDMovimiento: movID,
			ProductoID:   prodID,
			Tipo:         "INGRESO",
			Cantidad:     10,
			Motivo:       "Reposición inicial",
		}

		mov, err := queries.CreateMovimiento(ctx, paramsMov)
		if err != nil {
			t.Fatalf("Caso válido falló al crear movimiento de stock: %v", err)
		}
		if mov.Cantidad != 10 {
			t.Errorf("Esperaba cantidad 10, obtuve %d", mov.Cantidad)
		}
	})

	t.Run("Caso Inválido - Movimiento con Producto Inexistente (Falla de FK)", func(t *testing.T) {
		paramsMovInvalido := CreateMovimientoParams{
			IDMovimiento: 501,
			ProductoID:   99999, // Producto que no existe
			Tipo:         "EGRESO",
			Cantidad:     5,
			Motivo:       "Venta fantasma",
		}

		_, err := queries.CreateMovimiento(ctx, paramsMovInvalido)
		if err == nil {
			t.Errorf("Se esperaba un error de llave foránea al asociar un movimiento a un producto inexistente")
		}
	})

	// =========================================================================
	// 5. LIMPIEZA / DELETE
	// =========================================================================
	t.Run("Validar Eliminación de Movimiento y Producto", func(t *testing.T) {
		err := queries.DeleteMovimiento(ctx, movID)
		if err != nil {
			t.Fatalf("Error al eliminar movimiento: %v", err)
		}

		err = queries.DeleteProducto(ctx, prodID)
		if err != nil {
			t.Fatalf("Error al eliminar producto: %v", err)
		}
	})
}
