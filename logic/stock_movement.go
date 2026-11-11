package logic

import (
	"context"
	"errors"
	"fmt"

	"ctrlstock/db"
)

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

// CrearMovimiento centraliza la validación y el registro en la base de datos.
func CrearMovimiento(ctx context.Context, queries *db.Queries, params db.CreateMovimientoParams) (db.MovimientoStock, error) {
	if err := ValidateMovimiento(params.ProductoID, params.Tipo, params.Cantidad, params.Motivo); err != nil {
		return db.MovimientoStock{}, err
	}

	movimiento, err := queries.CreateMovimiento(ctx, params)
	if err != nil {
		return db.MovimientoStock{}, fmt.Errorf("error al registrar el movimiento de stock: %w", err)
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

// EliminarMovimiento encapsula el borrado de un registro en la bitácora.
func EliminarMovimiento(ctx context.Context, queries *db.Queries, id int32) error {
	err := queries.DeleteMovimiento(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el movimiento: %w", err)
	}
	return nil
}
