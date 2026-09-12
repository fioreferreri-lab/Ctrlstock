-- ============================================================================
-- CATEGORIA
-- ============================================================================

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
INSERT INTO categoria (id_categoria, nombre, padre_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCategoria :exec
UPDATE categoria
SET nombre = $1, padre_id = $2
WHERE id_categoria = $3;

-- name: DeleteCategoria :exec
DELETE FROM categoria 
WHERE id_categoria = $1;


-- ============================================================================
-- PRODUCTO
-- ============================================================================

-- name: GetProductoByID :one
SELECT * FROM producto 
WHERE id_producto = $1;

-- name: ListProductos :many
SELECT * FROM producto 
ORDER BY nombre;

-- name: CreateProducto :one
INSERT INTO producto (id_producto, nombre, descripcion, precio, stock, id_categoria)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProducto :exec
UPDATE producto
SET nombre = $1, descripcion = $2, precio = $3, stock = $4, id_categoria = $5
WHERE id_producto = $6;

-- name: DeleteProducto :exec
DELETE FROM producto 
WHERE id_producto = $1;


-- ============================================================================
-- MOVIMIENTO_STOCK
-- ============================================================================

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
INSERT INTO movimiento_stock (id_movimiento, producto_id, tipo, cantidad, motivo)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteMovimiento :exec
DELETE FROM movimiento_stock 
WHERE id_movimiento = $1;