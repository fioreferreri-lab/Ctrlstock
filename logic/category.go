package logic

import (
	"errors"
	"unicode/utf8"
)

// ValidateCategoria recibe el nombre ya limpiado (sin espacios al inicio/fin).
func ValidateCategoria(nombre string) error {
	if nombre == "" {
		return errors.New("el nombre de la categoría no puede estar vacío")
	}
	if utf8.RuneCountInString(nombre) > 100 {
		return errors.New("el nombre no puede superar los 100 caracteres")
	}
	return nil
}

// ValidatePadre se usa en el PUT: una categoría no puede ser su propio padre.
func ValidatePadre(id int32, padreID *int32) error {
	if padreID != nil && *padreID == id {
		return errors.New("una categoría no puede ser su propia categoría padre")
	}
	return nil
}
