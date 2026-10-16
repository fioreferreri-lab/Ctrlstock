-- 1. Borramos las tablas en orden inverso a sus dependencias
DROP TABLE IF EXISTS movimiento_stock;
DROP TABLE IF EXISTS producto;
DROP TABLE IF EXISTS categoria;

-- 2. Creación de tablas
CREATE TABLE categoria (
    id_categoria SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    padre_id INT NULL,
    CONSTRAINT CATEGORIA_Padre_fk FOREIGN KEY (padre_id) REFERENCES categoria(id_categoria) ON DELETE SET NULL
);

CREATE TABLE producto (
    id_producto SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255) NULL,
    precio DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    id_categoria INT NULL,
    CONSTRAINT PRODUCTO_CATEGORIA_fk FOREIGN KEY (id_categoria) REFERENCES categoria(id_categoria) ON DELETE SET NULL
);

CREATE TABLE movimiento_stock (
    id_movimiento SERIAL PRIMARY KEY,
    producto_id INT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    cantidad INT NOT NULL,
    motivo VARCHAR(100) NOT NULL,
    fecha TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT MOVIMIENTO_STOCK_PRODUCTO_fk FOREIGN KEY (producto_id) REFERENCES producto(id_producto) ON DELETE CASCADE
);