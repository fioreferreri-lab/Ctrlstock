package logic

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Category representa la tabla 'categoria' en Go con sus respectivos tags JSON.
type Category struct {
	IDCat   int    `json:"id_categoria"`
	Nombre  string `json:"nombre"`
	PadreID *int   `json:"padre_id,omitempty"` // *int permite valores NULL / nil de SQL
}

// ValidateCategory aplica las reglas de negocio básicas para la categoría.
func ValidateCategory(c Category) error {
	//se encarga de eliminar todos los espacios en blanco que estén al inicio y al final de una cadena de texto (string).
	nombreLimpio := strings.TrimSpace(c.Nombre)

	if nombreLimpio == "" {
		return errors.New("el nombre de la categoría no puede estar vacío")
	}

	if utf8.RuneCountInString(nombreLimpio) > 30 {
		return errors.New("el nombre de la categoría no puede superar los 30 caracteres")
	}

	// una categoría no puede ser su propio padre si manejas una ID asignada
	if c.PadreID != nil && c.IDCat > 0 && *c.PadreID == c.IDCat {
		return errors.New("una categoría no puede ser su propia categoría padre")
	}

	return nil
}
