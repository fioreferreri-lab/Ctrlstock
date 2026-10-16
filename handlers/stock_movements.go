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

type MovimientoInput struct {
	ProductoID int32  `json:"producto_id"`
	Tipo       string `json:"tipo"`
	Cantidad   int32  `json:"cantidad"`
	Motivo     string `json:"motivo"`
}

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

	if strings.Contains(path, "/") {
		http.Error(w, "Ruta no encontrada", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de movimiento inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getMovimiento(w, r, int32(id))
	case http.MethodDelete:
		s.deleteMovimiento(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listMovimientos(w http.ResponseWriter, r *http.Request) {
	lista, err := s.Queries.ListMovimientos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if lista == nil {
		lista = []db.MovimientoStock{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func (s *Server) createMovimiento(w http.ResponseWriter, r *http.Request) {
	var in MovimientoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	tipo := strings.ToUpper(strings.TrimSpace(in.Tipo))
	motivo := strings.TrimSpace(in.Motivo)
	if err := logic.ValidateMovimiento(in.ProductoID, tipo, in.Cantidad, motivo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := s.Queries.CreateMovimiento(r.Context(), db.CreateMovimientoParams{
		ProductoID: in.ProductoID,
		Tipo:       tipo,
		Cantidad:   in.Cantidad,
		Motivo:     motivo,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func (s *Server) getMovimiento(w http.ResponseWriter, r *http.Request, id int32) {
	m, err := s.Queries.GetMovimientoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Movimiento no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (s *Server) deleteMovimiento(w http.ResponseWriter, r *http.Request, id int32) {
	if err := s.Queries.DeleteMovimiento(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
