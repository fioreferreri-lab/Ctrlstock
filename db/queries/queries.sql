-- name: GetProductoByID :one
SELECT * FROM productos WHERE id = $1;

-- name: ListProductos :many
SELECT * FROM productos ORDER BY nombre;

-- name: CreateProducto :one
INSERT INTO productos (nombre, descripcion, precio, stock)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProducto :exec
UPDATE productos
SET nombre = $1, descripcion = $2, precio = $3, stock = $4
WHERE id = $5;

-- name: DeleteProducto :exec
DELETE FROM productos WHERE id = $1;