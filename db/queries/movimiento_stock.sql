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

-- name: RankingVentas :many
-- Productos ordenados por unidades vendidas desde una fecha.
-- Solo cuenta egresos con motivo "Venta" (no roturas ni ajustes).
SELECT p.id_producto, p.nombre, p.stock, c.nombre AS categoria,
       SUM(m.cantidad)::bigint AS vendidas
FROM movimiento_stock m
JOIN producto p ON p.id_producto = m.producto_id
JOIN categoria c ON c.id_categoria = p.id_categoria
WHERE m.tipo = 'EGRESO'
  AND m.motivo = 'Venta'
  AND m.fecha >= sqlc.arg(desde)::timestamp
GROUP BY p.id_producto, p.nombre, p.stock, c.nombre
ORDER BY vendidas DESC, p.nombre;
