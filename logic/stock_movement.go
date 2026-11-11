package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ctrlstock/db"
)

// ErrProductoInexistente se devuelve cuando el movimiento apunta a un producto que no existe.
var ErrProductoInexistente = errors.New("el producto seleccionado no existe")

// ValidateMovimiento recibe tipo (en mayúsculas) y motivo ya limpiados.
func ValidateMovimiento(productoID int32, tipo string, cantidad int32, motivo string) error {
	if productoID <= 0 {
		return errors.New("el ID de producto es obligatorio")
	}
	if tipo != "INGRESO" && tipo != "EGRESO" {
		return errors.New("el tipo de movimiento debe ser 'INGRESO' o 'EGRESO'")
	}
	if cantidad <= 0 {
		return errors.New("la cantidad del movimiento debe ser mayor a cero")
	}
	if motivo == "" {
		return errors.New("el motivo no puede estar vacío")
	}
	return nil
}

// CrearMovimiento valida el movimiento y, en UNA transacción:
//  1. suma (INGRESO) o resta (EGRESO) las unidades al stock del producto, y
//  2. registra el movimiento en el historial.
//
// O se guardan las dos cosas, o ninguna: nunca queda un cambio de stock sin
// registrar, ni un movimiento registrado que no haya cambiado el stock.
func CrearMovimiento(ctx context.Context, conn *sql.DB, queries *db.Queries, params db.CreateMovimientoParams) (db.MovimientoStock, error) {
	// 1. Validaciones de negocio (antes de tocar la base)
	if err := ValidateMovimiento(params.ProductoID, params.Tipo, params.Cantidad, params.Motivo); err != nil {
		return db.MovimientoStock{}, err
	}

	// 2. Abrimos la transacción. Si algo falla antes del Commit,
	//    el Rollback deshace todo lo que se hizo dentro de ella.
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return db.MovimientoStock{}, fmt.Errorf("error al iniciar la transacción: %w", err)
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx) // las mismas consultas de sqlc, pero dentro de la transacción

	// 3. Verificamos que el producto exista (y leemos su stock actual)
	producto, err := qtx.GetProductoByID(ctx, params.ProductoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.MovimientoStock{}, ErrProductoInexistente
		}
		return db.MovimientoStock{}, fmt.Errorf("error al buscar el producto: %w", err)
	}

	// 4. Actualizamos el stock: un egreso resta unidades
	delta := params.Cantidad
	if params.Tipo == "EGRESO" {
		delta = -params.Cantidad
	}
	_, err = qtx.AjustarStock(ctx, db.AjustarStockParams{
		Delta:      delta,
		IDProducto: params.ProductoID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// La consulta no actualiza si el stock quedaría negativo
			return db.MovimientoStock{}, fmt.Errorf(
				"stock insuficiente: hay %d unidades de %s y querés retirar %d",
				producto.Stock, producto.Nombre, params.Cantidad)
		}
		return db.MovimientoStock{}, fmt.Errorf("error al actualizar el stock: %w", err)
	}

	// 5. Registramos el movimiento en el historial
	movimiento, err := qtx.CreateMovimiento(ctx, params)
	if err != nil {
		return db.MovimientoStock{}, fmt.Errorf("error al registrar el movimiento de stock: %w", err)
	}

	// 6. Todo salió bien: confirmamos los dos cambios juntos
	if err := tx.Commit(); err != nil {
		return db.MovimientoStock{}, fmt.Errorf("error al confirmar el movimiento: %w", err)
	}

	return movimiento, nil
}

// ObtenerMovimiento encapsula la lectura de un movimiento específico por ID.
func ObtenerMovimiento(ctx context.Context, queries *db.Queries, id int32) (db.MovimientoStock, error) {
	movimiento, err := queries.GetMovimientoByID(ctx, id)
	if err != nil {
		return db.MovimientoStock{}, fmt.Errorf("error al buscar el movimiento: %w", err)
	}
	return movimiento, nil
}

// ListarMovimientos encapsula la obtención del historial completo de stock.
func ListarMovimientos(ctx context.Context, queries *db.Queries) ([]db.MovimientoStock, error) {
	lista, err := queries.ListMovimientos(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al listar los movimientos: %w", err)
	}
	if lista == nil {
		return []db.MovimientoStock{}, nil
	}
	return lista, nil
}

// ListarMovimientosByProducto filtra el historial para auditar un producto específico.
func ListarMovimientosByProducto(ctx context.Context, queries *db.Queries, productoID int32) ([]db.MovimientoStock, error) {
	lista, err := queries.ListMovimientosByProducto(ctx, productoID)
	if err != nil {
		return nil, fmt.Errorf("error al listar los movimientos del producto: %w", err)
	}
	if lista == nil {
		return []db.MovimientoStock{}, nil
	}
	return lista, nil
}

// Los movimientos NO se eliminan: son el historial de stock (registro de auditoría).
// Si hubo un error de carga, se corrige registrando el movimiento contrario
// (por ejemplo, un egreso que compense un ingreso equivocado).
