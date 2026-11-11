package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/db"
	"ctrlstock/logic"
	"ctrlstock/views"
)

// MovimientosHandler rige las rutas de /movimiento_stock
func (s *Server) MovimientosHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/movimiento_stock")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			s.listMovimientos(w, r)
		case http.MethodPost:
			s.createMovimiento(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		http.Error(w, "ID de movimiento inválido", http.StatusBadRequest)
		return
	}
	id32 := int32(id)

	switch r.Method {
	case http.MethodDelete:
		s.deleteMovimiento(w, r, id32)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// listMovimientos obtiene el historial y renderiza la vista HTML con Templ
func (s *Server) listMovimientos(w http.ResponseWriter, r *http.Request) {
	movimientos, err := logic.ListarMovimientos(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al listar los movimientos de stock", http.StatusInternalServerError)
		return
	}

	productos, err := logic.ListarProductos(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al listar los productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.MovimientoList(movimientos, productos).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
	}
}

// createMovimiento procesa el formulario clásico para registrar un movimiento
func (s *Server) createMovimiento(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	// Chequeo básico: que los campos numéricos sean números
	productoID, err := leerEntero(r, "producto_id", "El producto")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cantidad, err := leerEntero(r, "cantidad", "La cantidad")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.CreateMovimientoParams{
		ProductoID: productoID,
		Tipo:       strings.TrimSpace(r.FormValue("tipo")),
		Cantidad:   cantidad,
		Motivo:     strings.TrimSpace(r.FormValue("motivo")),
	}

	_, err = logic.CrearMovimiento(r.Context(), s.Queries, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Patrón PRG: Redirección limpia tras registrar
	http.Redirect(w, r, "/movimiento_stock", http.StatusSeeOther)
}

// deleteMovimiento elimina un registro de movimiento por ID y redirige
func (s *Server) deleteMovimiento(w http.ResponseWriter, r *http.Request, id int32) {
	if err := logic.EliminarMovimiento(r.Context(), s.Queries, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/movimiento_stock", http.StatusSeeOther)
}
