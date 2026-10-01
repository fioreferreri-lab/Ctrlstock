package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/logic"
)

// Lista global de productos en memoria
var products = []logic.Product{
	{
		ID:          1,
		Name:        "Teclado Mecánico",
		Description: stringPtr("Teclado RGB switch blue"),
		Price:       49.99,
		Stock:       15,
		CategoryID:  intPtr(1),
	},
	{
		ID:          2,
		Name:        "Mouse Gamer",
		Description: stringPtr("Mouse óptico 16000 DPI"),
		Price:       29.99,
		Stock:       30,
		CategoryID:  intPtr(1),
	},
}

// Funciones auxiliares para crear punteros en los datos de prueba
func intPtro(i int) *int         { return &i }
func stringPtr(s string) *string { return &s }

// ProductsHandler maneja la colección (/productos) -> GET (listar) y POST (crear)
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ProductHandler maneja un recurso individual (/productos/{id}) -> GET, PUT, DELETE
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProduct(w, r, id)
	case http.MethodPut:
		updateProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Operaciones CRUD

// GET /productos - Listar todos los productos
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// POST /productos - Crear un nuevo producto
func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct logic.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if newProduct.ID <= 0 {
		http.Error(w, "id_producto debe ser mayor a cero", http.StatusBadRequest)
		return
	}

	if err := logic.ValidateProduct(newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que el ID no esté duplicado
	if _, err := findProductByID(newProduct.ID); err == nil {
		http.Error(w, "ya existe un producto con ese ID", http.StatusConflict)
		return
	}

	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// GET /productos/{id} - Obtener un producto por ID
func getProduct(w http.ResponseWriter, r *http.Request, id int) {
	product, err := findProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// PUT /productos/{id} - Actualizar un producto existente
func updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findProductIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var updatedProduct logic.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logic.ValidateProduct(updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preservar el ID de la URL
	updatedProduct.ID = id
	products[index] = updatedProduct

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}

// DELETE /productos/{id} - Eliminar un producto
func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findProductIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	products = append(products[:index], products[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// Funciones auxiliares de búsqueda

func findProductByID(id int) (*logic.Product, error) {
	for i := range products {
		if products[i].ID == id {
			return &products[i], nil
		}
	}
	return nil, errors.New("producto no encontrado")
}

func findProductIndexByID(id int) (int, error) {
	for i, p := range products {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("producto no encontrado")
}
