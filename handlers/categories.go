package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/logic"
)

// Lista global de categorías en memoria (Estilo Ejercicio 3: sin Mutex)
var categories = []logic.Category{
	{IDCat: 1, Nombre: "Electrónica", PadreID: nil},
	{IDCat: 2, Nombre: "Computadoras", PadreID: intPtr(1)},
}

// Función auxiliar para crear punteros a int en la lista de prueba inicial
func intPtr(i int) *int { return &i }

// CategoriesHandler maneja la colección (/categories) -> GET (listar) y POST (crear)
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCategories(w, r)
	case http.MethodPost:
		createCategory(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// CategoryHandler maneja un recurso individual (/categories/{id}) -> GET, PUT, DELETE
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL (ejemplo: "/categories/2")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getCategory(w, r, id)
	case http.MethodPut:
		updateCategory(w, r, id)
	case http.MethodDelete:
		deleteCategory(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Operaciones CRUD

// GET /categories - Listar todas las categorías
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// POST /categories - Crear una nueva categoría
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCat logic.Category
	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateCategory(newCat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCat.IDCat = len(categories) + 1
	categories = append(categories, newCat)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCat)
}

// GET /categories/{id} - Obtener una categoría específica por ID
func getCategory(w http.ResponseWriter, r *http.Request, id int) {
	cat, err := findCategoryByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cat)
}

// PUT /categories/{id} - Actualizar una categoría existente
func updateCategory(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findCategoryIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var updatedCat logic.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCat); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateCategory(updatedCat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preservar el ID de la URL y actualizar la posición en el slice
	updatedCat.IDCat = id
	categories[index] = updatedCat

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCat)
}

// DELETE /categories/{id} - Eliminar una categoría
func deleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findCategoryIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Eliminar elemento del slice rebanando
	categories = append(categories[:index], categories[index+1:]...)

	// Responder con 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// Funciones auxiliares de búsqueda

func findCategoryByID(id int) (*logic.Category, error) {
	for i := range categories {
		if categories[i].IDCat == id {
			return &categories[i], nil
		}
	}
	return nil, errors.New("categoría no encontrada")
}

func findCategoryIndexByID(id int) (int, error) {
	for i, c := range categories {
		if c.IDCat == id {
			return i, nil
		}
	}
	return -1, errors.New("categoría no encontrada")
}
