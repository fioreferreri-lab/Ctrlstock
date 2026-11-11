package handlers

import (
	"database/sql"
	"errors"
	"log"
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

// listCategorias muestra la página de categorías con el formulario vacío (GET /categoria).
func (s *Server) listCategorias(w http.ResponseWriter, r *http.Request) {
	s.renderCategorias(w, r, http.StatusOK, views.DatosCategoriaForm{})
}

// renderCategorias arma la página de categorías. Se usa al entrar (formulario vacío)
// y cuando el formulario tuvo un error (con los datos escritos y el mensaje).
func (s *Server) renderCategorias(w http.ResponseWriter, r *http.Request, status int, form views.DatosCategoriaForm) {
	lista, err := logic.ListarCategorias(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al listar las categorías", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := views.CategoriaIndexPage("Categorías", lista, form).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar categorías: %v", err)
	}
}

// createCategoria procesa el formulario HTML clásico para dar de alta una categoría
func (s *Server) createCategoria(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderCategorias(w, r, http.StatusBadRequest, views.DatosCategoriaForm{Error: "el formulario es inválido"})
		return
	}

	// Lo que escribió el usuario: si algo falla, se lo devolvemos para que corrija
	form := views.DatosCategoriaForm{
		Nombre:  r.FormValue("nombre"),
		PadreID: r.FormValue("padre_id"),
	}
	mostrarError := func(err error) {
		form.Error = err.Error()
		s.renderCategorias(w, r, http.StatusBadRequest, form)
	}

	// Chequeo básico: que la categoría padre, si viene, sea un número
	padreID, err := leerEnteroOpcional(r, "padre_id", "La categoría padre")
	if err != nil {
		mostrarError(err)
		return
	}

	params := db.CreateCategoriaParams{
		Nombre:  strings.TrimSpace(r.FormValue("nombre")),
		PadreID: toNullInt32(padreID), // Utiliza el helper centralizado en utils.go
	}

	if _, err := logic.CrearCategoria(r.Context(), s.Queries, params); err != nil {
		mostrarError(err)
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

	// Chequeo básico: que la categoría padre, si viene, sea un número
	padreID, err := leerEnteroOpcional(r, "padre_id", "La categoría padre")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateCategoriaParams{
		IDCategoria: id,
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		PadreID:     toNullInt32(padreID),
	}

	err = logic.ActualizarCategoria(r.Context(), s.Queries, params)
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
