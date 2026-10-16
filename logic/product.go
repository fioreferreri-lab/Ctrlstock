package logic

import "errors"

// ValidateProducto recibe el nombre ya limpiado (sin espacios al inicio/fin).
func ValidateProducto(nombre string, precio float64, stock int32) error {
	if nombre == "" {
		return errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	if stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}
	return nil
}
