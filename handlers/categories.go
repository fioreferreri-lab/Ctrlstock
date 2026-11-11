package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/db"
	"ctrlstock/logic"
	"ctrlstock/views" // Importa plantillas de Templ
)

// CategoriasHandler rige las rutas de /categoria
func (s *Server) CategoriasHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/categoria")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			s.listCategorias(w, r)
		case http.MethodPost:
			s.createCategoria(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Manejo por ID para acciones específicas (/categoria/{id})
	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}
	id32 := int32(id)

	switch r.Method {
	case http.MethodPost, http.MethodPut:
		s.updateCategoria(w, r, id32)
	case http.MethodDelete:
		s.deleteCategoria(w, r, id32)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// listCategorias obtiene las categorías y renderiza la vista HTML con Templ
func (s *Server) listCategorias(w http.ResponseWriter, r *http.Request) {
	lista, err := logic.ListarCategorias(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al listar las categorías", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := views.CategoriaIndexPage("Categorías", lista).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
	}
}

// createCategoria procesa el formulario HTML clásico para dar de alta una categoría
func (s *Server) createCategoria(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	// Manejo seguro del campo opcional padre_id
	var padreID *int32
	padreStr := strings.TrimSpace(r.FormValue("padre_id"))
	if padreStr != "" {
		if pID, err := strconv.Atoi(padreStr); err == nil {
			p32 := int32(pID)
			padreID = &p32
		}
	}

	params := db.CreateCategoriaParams{
		Nombre:  strings.TrimSpace(r.FormValue("nombre")),
		PadreID: toNullInt32(padreID), // Utiliza el helper centralizado en utils.go
	}

	_, err := logic.CrearCategoria(r.Context(), s.Queries, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Patrón PRG: Redirección limpia tras guardar con éxito
	http.Redirect(w, r, "/categoria", http.StatusSeeOther)
}

// updateCategoria procesa la actualización de una categoría existente
func (s *Server) updateCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	var padreID *int32
	padreStr := strings.TrimSpace(r.FormValue("padre_id"))
	if padreStr != "" {
		if pID, err := strconv.Atoi(padreStr); err == nil {
			p32 := int32(pID)
			padreID = &p32
		}
	}

	params := db.UpdateCategoriaParams{
		IDCategoria: id,
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		PadreID:     toNullInt32(padreID),
	}

	err := logic.ActualizarCategoria(r.Context(), s.Queries, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/categoria", http.StatusSeeOther)
}

// deleteCategoria elimina una categoría por ID y redirige
func (s *Server) deleteCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	if err := logic.EliminarCategoria(r.Context(), s.Queries, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categoria", http.StatusSeeOther)
}
