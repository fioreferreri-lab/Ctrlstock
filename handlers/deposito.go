package handlers

import (
	"log"
	"net/http"

	"ctrlstock/logic"
	"ctrlstock/views"
)

// DepositoHandler muestra el árbol de categorías y productos con sus totales (GET /deposito).
func (s *Server) DepositoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	categorias, err := logic.ListarCategorias(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al obtener las categorías", http.StatusInternalServerError)
		return
	}

	productos, err := logic.ListarProductos(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al obtener los productos", http.StatusInternalServerError)
		return
	}

	// La lógica arma el árbol y calcula las inversiones; la vista solo lo muestra
	deposito := logic.ArmarDeposito(categorias, productos)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.DepositoPage("Depósito", deposito).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar el depósito: %v", err)
	}
}
