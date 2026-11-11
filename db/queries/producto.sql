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
-- El stock NO se edita acá: solo cambia con movimientos (ingresos y egresos).
UPDATE producto
SET nombre = $1, descripcion = $2, precio = $3, id_categoria = $4
WHERE id_producto = $5
RETURNING *;

-- name: DeleteProducto :exec
DELETE FROM producto
WHERE id_producto = $1;

-- name: AjustarStock :one
-- Suma unidades al stock (o resta, si delta es negativo).
-- Si el stock quedaría negativo, no modifica nada y no devuelve filas.
UPDATE producto
SET stock = stock + sqlc.arg(delta)::int
WHERE id_producto = sqlc.arg(id_producto)
  AND stock + sqlc.arg(delta)::int >= 0
RETURNING *;
