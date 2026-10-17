-- name: GetProductoByID :one
SELECT * FROM producto
WHERE id_producto = $1;

-- name: ListProductos :many
SELECT * FROM producto
ORDER BY nombre;

-- name: CreateProducto :one
INSERT INTO producto (nombre, descripcion, precio, stock, id_categoria)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProducto :one
UPDATE producto
SET nombre = $1, descripcion = $2, precio = $3, stock = $4, id_categoria = $5
WHERE id_producto = $6
RETURNING *;

-- name: DeleteProducto :exec
DELETE FROM producto
WHERE id_producto = $1;
