package logic

import "errors"

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
