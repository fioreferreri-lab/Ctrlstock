package logic

import (
	"context"
	"ctrlstock/db"
	"errors"
	"fmt"
	"strings"
)

// ValidateProducto aplica reglas de negocio a un producto.
func ValidateProducto(nombre string, precio float64, stock int32) error {
	if strings.TrimSpace(nombre) == "" {
		return errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	if stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}
	return nil
}

// CrearProducto centraliza la validación y la inserción en la base de datos.
func CrearProducto(ctx context.Context, queries *db.Queries, params db.CreateProductoParams) (db.Producto, error) {
	// Parseamos el string del precio (db.CreateProductoParams usa string para DECIMAL) a float64 para validarlo
	var precioFloat float64
	fmt.Sscanf(params.Precio, "%f", &precioFloat)

	// 1. Validaciones importantes (Regla de negocio)
	if err := ValidateProducto(params.Nombre, precioFloat, params.Stock); err != nil {
		return db.Producto{}, err
	}

	// 2. Consulta a la base de datos (Persistencia delegada a sqlc)
	producto, err := queries.CreateProducto(ctx, params)
	if err != nil {
		return db.Producto{}, fmt.Errorf("error al guardar el producto en la base de datos: %w", err)
	}

	return producto, nil
}

// ObtenerProducto encapsula la lectura de un producto por ID.
func ObtenerProducto(ctx context.Context, queries *db.Queries, id int32) (db.Producto, error) {
	producto, err := queries.GetProductoByID(ctx, id)
	if err != nil {
		return db.Producto{}, fmt.Errorf("error al buscar el producto: %w", err)
	}
	return producto, nil
}

// ListarProductos encapsula la obtención del catálogo completo.
func ListarProductos(ctx context.Context, queries *db.Queries) ([]db.Producto, error) {
	lista, err := queries.ListProductos(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al listar los productos: %w", err)
	}
	if lista == nil {
		return []db.Producto{}, nil
	}
	return lista, nil
}

// ActualizarProducto centraliza la validación y la actualización en la base de datos.
func ActualizarProducto(ctx context.Context, queries *db.Queries, params db.UpdateProductoParams) (db.Producto, error) {
	var precioFloat float64
	fmt.Sscanf(params.Precio, "%f", &precioFloat)

	// 1. Validaciones importantes (Regla de negocio)
	if err := ValidateProducto(params.Nombre, precioFloat, params.Stock); err != nil {
		return db.Producto{}, err
	}

	// 2. Actualización en la base de datos
	producto, err := queries.UpdateProducto(ctx, params)
	if err != nil {
		return db.Producto{}, fmt.Errorf("error al actualizar el producto: %w", err)
	}

	return producto, nil
}

// EliminarProducto encapsula el borrado de un producto.
func EliminarProducto(ctx context.Context, queries *db.Queries, id int32) error {
	err := queries.DeleteProducto(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el producto: %w", err)
	}
	return nil
}
