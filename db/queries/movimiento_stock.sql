-- name: GetMovimientoByID :one
SELECT * FROM movimiento_stock
WHERE id_movimiento = $1;

-- name: ListMovimientos :many
SELECT * FROM movimiento_stock
ORDER BY fecha DESC;

-- name: ListMovimientosByProducto :many
SELECT * FROM movimiento_stock
WHERE producto_id = $1
ORDER BY fecha DESC;

-- name: CreateMovimiento :one
INSERT INTO movimiento_stock (producto_id, tipo, cantidad, motivo)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteMovimiento :exec
DELETE FROM movimiento_stock
WHERE id_movimiento = $1;
