package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/logic"
)

// Lista global de movimientos en memoria
var movimientos = []logic.Stockmovement{
	{IDMovimiento: 1, Tipo: "INGRESO", Cantidad: 50},
	{IDMovimiento: 2, Tipo: "EGRESO", Cantidad: 10},
}

// MovimientosHandler maneja la colección (/movimientos) -> GET (listar) y POST (crear)
func MovimientosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMovimientos(w, r)
	case http.MethodPost:
		createMovimiento(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// MovimientoHandler maneja un recurso individual (/movimientos/{id}) -> GET, PUT, DELETE
func MovimientoHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL (ejemplo: "/movimientos/2")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "ID de movimiento inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getMovimiento(w, r, id)
	case http.MethodPut:
		updateMovimiento(w, r, id)
	case http.MethodDelete:
		deleteMovimiento(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Operaciones CRUD

// GET /movimientos - Listar todos los movimientos
func getMovimientos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movimientos)
}

// POST /movimientos - Crear un nuevo movimiento
func createMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMov logic.Stockmovement
	if err := json.NewDecoder(r.Body).Decode(&newMov); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateMovimiento(newMov); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMov.IDMovimiento = len(movimientos) + 1
	movimientos = append(movimientos, newMov)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMov)
}

// GET /movimientos/{id} - Obtener un movimiento específico por ID
func getMovimiento(w http.ResponseWriter, r *http.Request, id int) {
	mov, err := findMovimientoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mov)
}

// PUT /movimientos/{id} - Actualizar un movimiento existente
func updateMovimiento(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findMovimientoIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var updatedMov logic.Stockmovement
	if err := json.NewDecoder(r.Body).Decode(&updatedMov); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateMovimiento(updatedMov); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preservar el ID de la URL y actualizar la posición en el slice
	updatedMov.IDMovimiento = id
	movimientos[index] = updatedMov

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMov)
}

// DELETE /movimientos/{id} - Eliminar un movimiento
func deleteMovimiento(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findMovimientoIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Eliminar elemento del slice rebanando
	movimientos = append(movimientos[:index], movimientos[index+1:]...)

	// Responder con 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// Funciones auxiliares de búsqueda

func findMovimientoByID(id int) (*logic.Stockmovement, error) {
	for i := range movimientos {
		if movimientos[i].IDMovimiento == id {
			return &movimientos[i], nil
		}
	}
	return nil, errors.New("movimiento no encontrado")
}

func findMovimientoIndexByID(id int) (int, error) {
	for i, m := range movimientos {
		if m.IDMovimiento == id {
			return i, nil
		}
	}
	return -1, errors.New("movimiento no encontrado")
}
