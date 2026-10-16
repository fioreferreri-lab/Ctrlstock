package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "ctrlstock/db"
	"ctrlstock/logic"
)

type ProductoInput struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
	Precio      float64 `json:"precio"`
	Stock       int32   `json:"stock"`
	IDCategoria *int32  `json:"id_categoria,omitempty"`
}

func (s *Server) ProductosHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/producto")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			s.listProductos(w, r)
		case http.MethodPost:
			s.createProducto(w, r)
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
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getProducto(w, r, int32(id))
	case http.MethodPut:
		s.updateProducto(w, r, int32(id))
	case http.MethodDelete:
		s.deleteProducto(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listProductos(w http.ResponseWriter, r *http.Request) {
	lista, err := s.Queries.ListProductos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if lista == nil {
		lista = []db.Producto{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func (s *Server) createProducto(w http.ResponseWriter, r *http.Request) {
	var in ProductoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(in.Nombre)
	if err := logic.ValidateProducto(nombre, in.Precio, in.Stock); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.Queries.CreateProducto(r.Context(), db.CreateProductoParams{
		Nombre:      nombre,
		Descripcion: toNullString(in.Descripcion),
		Precio:      fmt.Sprintf("%.2f", in.Precio),
		Stock:       in.Stock,
		IDCategoria: toNullInt32(in.IDCategoria),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (s *Server) getProducto(w http.ResponseWriter, r *http.Request, id int32) {
	p, err := s.Queries.GetProductoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) updateProducto(w http.ResponseWriter, r *http.Request, id int32) {
	var in ProductoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(in.Nombre)
	if err := logic.ValidateProducto(nombre, in.Precio, in.Stock); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.Queries.UpdateProducto(r.Context(), db.UpdateProductoParams{
		IDProducto:  id,
		Nombre:      nombre,
		Descripcion: toNullString(in.Descripcion),
		Precio:      fmt.Sprintf("%.2f", in.Precio),
		Stock:       in.Stock,
		IDCategoria: toNullInt32(in.IDCategoria),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) deleteProducto(w http.ResponseWriter, r *http.Request, id int32) {
	if err := s.Queries.DeleteProducto(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
