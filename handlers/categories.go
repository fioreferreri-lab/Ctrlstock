package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "ctrlstock/db"
	"ctrlstock/logic"
)

type CategoriaInput struct {
	Nombre  string `json:"nombre"`
	PadreID *int32 `json:"padre_id,omitempty"`
}

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

	if strings.Contains(path, "/") {
		http.Error(w, "Ruta no encontrada", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getCategoria(w, r, int32(id))
	case http.MethodPut:
		s.updateCategoria(w, r, int32(id))
	case http.MethodDelete:
		s.deleteCategoria(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listCategorias(w http.ResponseWriter, r *http.Request) {
	lista, err := s.Queries.ListCategorias(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if lista == nil {
		lista = []db.Categorium{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func (s *Server) createCategoria(w http.ResponseWriter, r *http.Request) {
	var in CategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(in.Nombre)
	if err := logic.ValidateCategoria(nombre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := s.Queries.CreateCategoria(r.Context(), db.CreateCategoriaParams{
		Nombre:  nombre,
		PadreID: toNullInt32(in.PadreID),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (s *Server) getCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	c, err := s.Queries.GetCategoriaByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (s *Server) updateCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	var in CategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(in.Nombre)
	if err := logic.ValidateCategoria(nombre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := logic.ValidatePadre(id, in.PadreID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// La query usa ":one" con "RETURNING *": si el id no existe, sqlc devuelve
	// sql.ErrNoRows, por lo que no hace falta un Get previo.
	c, err := s.Queries.UpdateCategoria(r.Context(), db.UpdateCategoriaParams{
		IDCategoria: id,
		Nombre:      nombre,
		PadreID:     toNullInt32(in.PadreID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (s *Server) deleteCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	if err := s.Queries.DeleteCategoria(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
