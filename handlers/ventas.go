package handlers

import (
	"log"
	"net/http"
	"strconv"

	"ctrlstock/logic"
	"ctrlstock/views"
)

// VentasHandler muestra el ranking de productos más vendidos (GET /ventas?dias=30).
func (s *Server) VentasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Chequeo básico: si "dias" no es un número, se usa el período por defecto.
	// Qué períodos son válidos lo decide la capa de lógica.
	dias, err := strconv.Atoi(r.URL.Query().Get("dias"))
	if err != nil {
		dias = 30
	}

	ranking, err := logic.ObtenerRankingVentas(r.Context(), s.Queries, dias)
	if err != nil {
		http.Error(w, "Error al obtener el ranking de ventas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.VentasPage("Más vendidos", ranking).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar el ranking: %v", err)
	}
}
