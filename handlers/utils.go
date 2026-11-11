package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "ctrlstock/db"
)

// Server agrupa las dependencias que necesitan los handlers.
type Server struct {
	Queries *db.Queries
}

func toNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{String: s, Valid: s != ""}
}

func toNullInt32(n *int32) sql.NullInt32 {
	if n == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *n, Valid: true}
}

// ---------------------------------------------------------------
// Chequeo básico de los datos del formulario (responsabilidad del handler):
// solo verifican que el texto recibido tenga el FORMATO correcto
// (que sea un número). Las reglas de negocio (precio > 0, stock >= 0,
// nombre no vacío, etc.) las valida la capa de lógica.
// ---------------------------------------------------------------

// leerEntero lee un campo numérico entero obligatorio.
func leerEntero(r *http.Request, campo string, etiqueta string) (int32, error) {
	texto := strings.TrimSpace(r.FormValue(campo))
	n, err := strconv.ParseInt(texto, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s debe ser un número entero", etiqueta)
	}
	return int32(n), nil
}

// leerEnteroOpcional lee un campo entero que puede venir vacío (devuelve nil).
func leerEnteroOpcional(r *http.Request, campo string, etiqueta string) (*int32, error) {
	if strings.TrimSpace(r.FormValue(campo)) == "" {
		return nil, nil
	}
	n, err := leerEntero(r, campo, etiqueta)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// leerDecimal lee un campo numérico con decimales (ej. un precio).
func leerDecimal(r *http.Request, campo string, etiqueta string) (float64, error) {
	texto := strings.TrimSpace(r.FormValue(campo))
	v, err := strconv.ParseFloat(texto, 64)
	if err != nil {
		return 0, fmt.Errorf("%s debe ser un número", etiqueta)
	}
	return v, nil
}
