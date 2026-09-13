-- 1. Primero borramos las tablas en orden inverso a sus dependencias
DROP TABLE IF EXISTS movimiento_stock;
DROP TABLE IF EXISTS producto;
DROP TABLE IF EXISTS categoria;

-- 2. Ahora sí, las creamos en orden normal
CREATE TABLE categoria (
    id_categoria int NOT NULL,
    nombre varchar(100) NOT NULL,
    padre_id int NULL,
    CONSTRAINT CATEGORIA_pk PRIMARY KEY (id_categoria),
    CONSTRAINT CATEGORIA_Padre_fk FOREIGN KEY (padre_id) REFERENCES categoria(id_categoria) ON DELETE SET NULL
);

CREATE TABLE producto (
    id_producto int NOT NULL,
    nombre varchar(100) NOT NULL,
    descripcion varchar(255) NULL,
    precio decimal(10,2) NOT NULL,
    stock int NOT NULL,
    id_categoria int NULL,
    CONSTRAINT PRODUCTO_pk PRIMARY KEY (id_producto),
    CONSTRAINT PRODUCTO_CATEGORIA_fk FOREIGN KEY (id_categoria) REFERENCES categoria(id_categoria)
);

CREATE TABLE movimiento_stock(
    id_movimiento int NOT NULL,
    producto_id int NOT NULL,
    tipo varchar(20) NOT NULL,
    cantidad int NOT NULL,
    motivo varchar(100) NOT NULL,
    fecha timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT MOVIMIENTO_STOCK_pk PRIMARY KEY (id_movimiento),
    CONSTRAINT MOVIMIENTO_STOCK_PRODUCTO_fk FOREIGN KEY (producto_id) REFERENCES producto(id_producto) ON DELETE CASCADE
);