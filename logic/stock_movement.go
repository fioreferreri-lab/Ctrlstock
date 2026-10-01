package logic

import (
	"errors"
	"strings"
)

// Movimiento representa la entidad 'movimiento' en Go con sus respectivos tags JSON.
type Stockmovement struct {
	IDMovimiento int    `json:"id_movimiento"`
	Tipo         string `json:"tipo"`
	Cantidad     int    `json:"cantidad"`
}

// ValidateMovimiento aplica las reglas de negocio básicas para el movimiento.
func ValidateMovimiento(m Stockmovement) error {
	tipoLimpio := strings.TrimSpace(strings.ToUpper(m.Tipo))

	if tipoLimpio == "" {
		return errors.New("el tipo de movimiento no puede estar vacío")
	}

	if tipoLimpio != "INGRESO" && tipoLimpio != "EGRESO" {
		return errors.New("el tipo de movimiento debe ser 'INGRESO' o 'EGRESO'")
	}

	if m.Cantidad <= 0 {
		return errors.New("la cantidad del movimiento debe ser mayor a cero")
	}

	return nil
}
