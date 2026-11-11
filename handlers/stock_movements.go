package handlers

import (
	"log"
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

// listMovimientos muestra la página de movimientos con el formulario vacío (GET /movimiento_stock).
func (s *Server) listMovimientos(w http.ResponseWriter, r *http.Request) {
	s.renderMovimientos(w, r, http.StatusOK, views.DatosMovimientoForm{})
}

// renderMovimientos arma la página de movimientos. Se usa al entrar (formulario vacío)
// y cuando el formulario tuvo un error (con los datos escritos y el mensaje).
func (s *Server) renderMovimientos(w http.ResponseWriter, r *http.Request, status int, form views.DatosMovimientoForm) {
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
	w.WriteHeader(status)
	if err := views.MovimientoIndexPage("Movimientos de stock", movimientos, productos, form).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar movimientos: %v", err)
	}
}

// createMovimiento procesa el formulario clásico para registrar un movimiento
func (s *Server) createMovimiento(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderMovimientos(w, r, http.StatusBadRequest, views.DatosMovimientoForm{Error: "el formulario es inválido"})
		return
	}

	// Lo que escribió el usuario: si algo falla, se lo devolvemos para que corrija
	form := views.DatosMovimientoForm{
		ProductoID: r.FormValue("producto_id"),
		Tipo:       r.FormValue("tipo"),
		Cantidad:   r.FormValue("cantidad"),
		Motivo:     r.FormValue("motivo"),
	}
	mostrarError := func(err error) {
		form.Error = err.Error()
		s.renderMovimientos(w, r, http.StatusBadRequest, form)
	}

	// Chequeo básico: que los campos numéricos sean números
	productoID, err := leerEntero(r, "producto_id", "El producto")
	if err != nil {
		mostrarError(err)
		return
	}
	cantidad, err := leerEntero(r, "cantidad", "La cantidad")
	if err != nil {
		mostrarError(err)
		return
	}

	params := db.CreateMovimientoParams{
		ProductoID: productoID,
		Tipo:       strings.TrimSpace(r.FormValue("tipo")),
		Cantidad:   cantidad,
		Motivo:     strings.TrimSpace(r.FormValue("motivo")),
	}

	if _, err := logic.CrearMovimiento(r.Context(), s.Queries, params); err != nil {
		mostrarError(err)
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
