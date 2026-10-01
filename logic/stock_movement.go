package logic

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// Movimiento representa la entidad 'movimiento_stock' en Go.
type Stockmovement struct {
	IDMovimiento int       `json:"id_movimiento"`
	ProductoID   int       `json:"producto_id"`
	Tipo         string    `json:"tipo"`
	Cantidad     int       `json:"cantidad"`
	Motivo       string    `json:"motivo"`
	Fecha        time.Time `json:"fecha,omitempty"`
}

// ValidateMovimiento aplica las reglas de validación básicas del movimiento.
func ValidateMovimiento(m Stockmovement) error {
	// 1. Validar ProductoID
	if m.ProductoID <= 0 {
		return errors.New("el ID de producto debe ser válido y mayor a cero")
	}

	// 2. Validar Tipo de movimiento
	tipoLimpio := strings.TrimSpace(strings.ToUpper(m.Tipo))
	if tipoLimpio == "" {
		return errors.New("el tipo de movimiento no puede estar vacío")
	}
	if tipoLimpio != "INGRESO" && tipoLimpio != "EGRESO" {
		return errors.New("el tipo de movimiento debe ser 'INGRESO' o 'EGRESO'")
	}

	// 3. Validar Cantidad
	if m.Cantidad <= 0 {
		return errors.New("la cantidad del movimiento debe ser mayor a cero")
	}

	// 4. Validar Motivo
	motivoLimpio := strings.TrimSpace(m.Motivo)
	if motivoLimpio == "" {
		return errors.New("el motivo del movimiento no puede estar vacío")
	}
	if utf8.RuneCountInString(motivoLimpio) > 100 {
		return errors.New("el motivo no puede superar los 100 caracteres")
	}

	return nil
}
