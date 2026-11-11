package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ctrlstock/db"
	"ctrlstock/logic"
	"ctrlstock/views" // Importamos plantillas de Templ
)

// ProductosHandler rige tanto la colección (/producto) como recursos individuales (/producto/{id})
func (s *Server) ProductosHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/producto")
	path = strings.Trim(path, "/")

	// Si la ruta está vacía (/producto), manejamos listar y crear
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

	// Si hay un ID en la ruta (/producto/{id}), manejamos acciones específicas por ID
	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}
	id32 := int32(id)

	switch r.Method {
	case http.MethodPost, http.MethodPut:
		s.updateProducto(w, r, id32)
	case http.MethodDelete:
		s.deleteProducto(w, r, id32)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// listProductos renderiza la vista HTML principal usando Templ
func (s *Server) listProductos(w http.ResponseWriter, r *http.Request) {
	// 1. Obtenemos datos desde la capa lógica
	lista, err := logic.ListarProductos(r.Context(), s.Queries)
	if err != nil {
		http.Error(w, "Error al obtener el inventario", http.StatusInternalServerError)
		return
	}

	// 2. TP5: Renderizamos el HTML con Templ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.IndexPage("Inventario de Productos", lista).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
	}
}

// createProducto procesa el formulario HTML clásico para dar de alta un producto
func (s *Server) createProducto(w http.ResponseWriter, r *http.Request) {
	// 1. Parseamos el formulario clásico enviado por el navegador
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	// Manejo seguro del ID de categoría
	idCatStr := strings.TrimSpace(r.FormValue("id_categoria"))
	var idCategoria *int32
	if idCatStr != "" {
		if id, err := strconv.Atoi(idCatStr); err == nil {
			id32 := int32(id)
			idCategoria = &id32
		}
	}

	// Mapeo de parámetros para sqlc
	params := db.CreateProductoParams{
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Descripcion: toNullString(r.FormValue("descripcion")),
		Precio:      fmt.Sprintf("%.2f", precio),
		Stock:       int32(stock),
		IDCategoria: toNullInt32(idCategoria),
	}

	// 2. Delegamos validación y guardado a la capa lógica
	_, err := logic.CrearProducto(r.Context(), s.Queries, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Patrón PRG (Post/Redirect/Get): Redirigimos a la lista para actualizar la vista
	http.Redirect(w, r, "/producto", http.StatusSeeOther)
}

// updateProducto procesa la actualización de un producto existente por ID
func (s *Server) updateProducto(w http.ResponseWriter, r *http.Request, id int32) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	idCatStr := strings.TrimSpace(r.FormValue("id_categoria"))
	var idCategoria *int32
	if idCatStr != "" {
		if idVal, err := strconv.Atoi(idCatStr); err == nil {
			id32 := int32(idVal)
			idCategoria = &id32
		}
	}

	params := db.UpdateProductoParams{
		IDProducto:  id,
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Descripcion: toNullString(r.FormValue("descripcion")),
		Precio:      fmt.Sprintf("%.2f", precio),
		Stock:       int32(stock),
		IDCategoria: toNullInt32(idCategoria),
	}

	// Delegamos a la capa lógica
	if _, err := logic.ActualizarProducto(r.Context(), s.Queries, params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/producto", http.StatusSeeOther)
}

// deleteProducto elimina un producto por su ID y redirige
func (s *Server) deleteProducto(w http.ResponseWriter, r *http.Request, id int32) {
	if err := logic.EliminarProducto(r.Context(), s.Queries, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/producto", http.StatusSeeOther)
}
