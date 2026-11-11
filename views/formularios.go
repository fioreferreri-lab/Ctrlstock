package views

import (
	"strings"
	"unicode"
)

// Los formularios guardan lo que escribió el usuario y el error (si hubo).
// Cuando un dato es inválido, el handler vuelve a mostrar la página con
// estos datos: así el usuario ve qué falló y no tiene que escribir todo de nuevo.

// DatosProductoForm son los valores del formulario de producto.
type DatosProductoForm struct {
	Nombre      string
	Descripcion string
	Precio      string
	Stock       string
	IDCategoria string
	Error       string
}

// DatosCategoriaForm son los valores del formulario de categoría.
type DatosCategoriaForm struct {
	Nombre  string
	PadreID string
	Error   string
}

// DatosMovimientoForm son los valores del formulario de movimiento.
type DatosMovimientoForm struct {
	ProductoID string
	Tipo       string
	Cantidad   string
	Motivo     string
	Error      string
}

// valorOPorDefecto devuelve el valor escrito o, si está vacío, uno por defecto.
func valorOPorDefecto(valor string, porDefecto string) string {
	if valor == "" {
		return porDefecto
	}
	return valor
}

// capitalizar pone en mayúscula la primera letra del mensaje de error.
func capitalizar(mensaje string) string {
	runas := []rune(strings.TrimSpace(mensaje))
	if len(runas) == 0 {
		return ""
	}
	runas[0] = unicode.ToUpper(runas[0])
	return string(runas)
}
