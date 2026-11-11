package logic

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/lib/pq"

	"ctrlstock/db"
)

// ErrCategoriaDuplicada se devuelve cuando ya existe una categoría con ese nombre.
var ErrCategoriaDuplicada = errors.New("ya existe una categoría con ese nombre")

// esNombreDuplicado detecta si la base rechazó el dato por la restricción UNIQUE.
// 23505 es el código de Postgres para "valor duplicado".
func esNombreDuplicado(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

// ValidateCategoria recibe el nombre ya limpiado (sin espacios al inicio/fin).
func ValidateCategoria(nombre string) error {
	if nombre == "" {
		return errors.New("el nombre de la categoría no puede estar vacío")
	}
	if utf8.RuneCountInString(nombre) > 100 {
		return errors.New("el nombre no puede superar los 100 caracteres")
	}
	return nil
}

// ValidatePadre se usa en el PUT: una categoría no puede ser su propio padre.
func ValidatePadre(id int32, padreID *int32) error {
	if padreID != nil && *padreID == id {
		return errors.New("una categoría no puede ser su propia categoría padre")
	}
	return nil
}

// CrearCategoria centraliza la validación y la inserción en la base de datos.
func CrearCategoria(ctx context.Context, queries *db.Queries, params db.CreateCategoriaParams) (db.Categoria, error) {
	if err := ValidateCategoria(params.Nombre); err != nil {
		return db.Categoria{}, err
	}

	categoria, err := queries.CreateCategoria(ctx, params)
	if err != nil {
		if esNombreDuplicado(err) {
			return db.Categoria{}, ErrCategoriaDuplicada
		}
		return db.Categoria{}, fmt.Errorf("error al guardar la categoría: %w", err)
	}

	return categoria, nil
}

// ObtenerCategoria encapsula la lectura de una categoría por ID.
func ObtenerCategoria(ctx context.Context, queries *db.Queries, id int32) (db.Categoria, error) {
	categoria, err := queries.GetCategoriaByID(ctx, id)
	if err != nil {
		return db.Categoria{}, fmt.Errorf("error al buscar la categoría: %w", err)
	}
	return categoria, nil
}

// ListarCategorias encapsula la obtención del catálogo completo.
func ListarCategorias(ctx context.Context, queries *db.Queries) ([]db.Categoria, error) {
	lista, err := queries.ListCategorias(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al listar las categorías: %w", err)
	}
	if lista == nil {
		return []db.Categoria{}, nil
	}
	return lista, nil
}

// ActualizarCategoria centraliza la validación y la actualización en la base de datos.
func ActualizarCategoria(ctx context.Context, queries *db.Queries, params db.UpdateCategoriaParams) error {
	if err := ValidateCategoria(params.Nombre); err != nil {
		return err
	}
	// Extraemos el valor del padre para la validación lógica
	var padreID *int32
	if params.PadreID.Valid {
		padreID = &params.PadreID.Int32
	}

	if err := ValidatePadre(params.IDCategoria, padreID); err != nil {
		return err
	}

	_, err := queries.UpdateCategoria(ctx, params)
	if err != nil {
		if esNombreDuplicado(err) {
			return ErrCategoriaDuplicada
		}
		return fmt.Errorf("error al actualizar la categoría: %w", err)
	}

	return nil
}

// EliminarCategoria encapsula el borrado de una categoría.
func EliminarCategoria(ctx context.Context, queries *db.Queries, id int32) error {
	err := queries.DeleteCategoria(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la categoría: %w", err)
	}
	return nil
}
