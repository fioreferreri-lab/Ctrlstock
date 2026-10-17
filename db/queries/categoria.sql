-- name: GetCategoriaByID :one
SELECT * FROM categoria
WHERE id_categoria = $1;

-- name: ListCategorias :many
SELECT * FROM categoria
ORDER BY nombre;

-- name: ListCategoriasRaiz :many
SELECT * FROM categoria
WHERE padre_id IS NULL
ORDER BY nombre;

-- name: ListSubcategorias :many
SELECT * FROM categoria
WHERE padre_id = $1
ORDER BY nombre;

-- name: CreateCategoria :one
INSERT INTO categoria (nombre, padre_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCategoria :one
UPDATE categoria
SET nombre = $1, padre_id = $2
WHERE id_categoria = $3
RETURNING *;

-- name: DeleteCategoria :exec
DELETE FROM categoria
WHERE id_categoria = $1;
