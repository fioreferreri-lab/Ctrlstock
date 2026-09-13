# CtrlStock: Documentación Técnica y de Arquitectura

**CtrlStock** es una aplicación web robusta y moderna para la gestión integral de stock e inventario, diseñada con un enfoque en la mantenibilidad, el tipado seguro y la reproducibilidad del entorno. Permite la administración eficiente de productos, jerarquía de categorías y el registro auditable de todos los movimientos de inventario.

---

## 1. Arquitectura y Tecnologías Utilizadas

El sistema está construido utilizando un stack tecnológico moderno, ligero y altamente performante:

*   **Lenguaje Principal:** **Go (Golang)**, elegido por su alta concurrencia, rendimiento excepcional y tipado estático robusto.
*   **Base de Datos:** **PostgreSQL**, motor relacional potente y confiable para la integridad transaccional de los datos de inventario.
*   **Contenedorización:** **Docker & Docker Compose**, que garantizan un entorno de base de datos aislado, reproducible y libre de dependencias locales complejas.
*   **Generación de Código:** **sqlc**, herramienta que compila consultas SQL nativas en código Go fuertemente tipado, eliminando la necesidad de ORMs pesados y previniendo errores en tiempo de ejecución.
*   **Automatización:** **GNU Make**, para estandarizar y simplificar las tareas de desarrollo, pruebas y despliegue local.
*   **Driver de Conexión:** `lib/pq`, driver oficial de PostgreSQL para Go que gestiona la comunicación con la base de datos.

---

## 2. Requisitos del Sistema

Para poder ejecutar la automatización, persistencia, compilación y testeo del proyecto, es necesario contar con las siguientes herramientas instaladas en el sistema operativo:

1.  **Go:** Versión `>= 1.22.2`.
2.  **Docker y Docker Compose:** Indispensables para levantar el contenedor con la instancia de PostgreSQL.
3.  **Make:** Utilizado para ejecutar las recetas de automatización de pruebas y flujos de trabajo.
4.  **SQLC:** Utilizado para la generación automática de código Go a partir de las consultas SQL definidas para PostgreSQL.

---

## 3. Estructura del Proyecto

La disposición de los directorios y archivos principales del repositorio sigue una estructura limpia y modular:

```text
SISTEMA_STOCK/
│
├── db/
│   ├── scheme/
│   │   └── schema.sql        # Define las tablas: categoria, producto y movimiento_stock
│   ├── queries/
│   │   └── queries.sql       # Consultas SQL optimizadas utilizadas por sqlc
│   ├── models.go             # Código Go generado automáticamente por sqlc (estructuras de datos)
│   └── queries.sql.go        # Código Go generado automáticamente por sqlc (métodos CRUD)
│
├── docker-compose.yml        # Configuración del contenedor de PostgreSQL (desarrollo y testing)
└── Makefile                  # Automatización de generación de código, DB y tests
```

---

## 4. Modelo de Datos

El diseño relacional de la base de datos se compone de tres entidades principales que garantizan la trazabilidad y organización del inventario:

### A. Categoría (`categoria`)
Organiza los productos en una jerarquía flexible, permitiendo la creación de subcategorías mediante una relación recursiva consigo misma (`padre_id`).
*   **Campos clave:** ID, nombre, descripción, `padre_id` (autoreferencia opcional).

### B. Producto (`producto`)
Representa cada ítem individual dentro del catálogo del e-commerce.
*   **Campos clave:** ID, nombre, descripción, precio unitario, stock actual, `categoria_id` (clave foránea).

### C. Movimiento de Stock (`movimiento_stock`)
Registra de forma histórica y auditable cada entrada, salida o ajuste de inventario (por ejemplo: compras a proveedores, ventas a clientes, mermas o ajustes manuales). Esto evita mantener únicamente el número final, asegurando una auditoría completa.
*   **Campos clave:** ID, `producto_id`, tipo de movimiento (entrada/salida), cantidad, motivo/descripción, timestamp.

---

## 5. Persistencia de Datos

El almacenamiento y la gestión de la persistencia se realizan mediante PostgreSQL ejecutado en contenedores Docker:

*   **Cadena de Conexión:** La aplicación se conecta utilizando el formato estándar:
    `postgres://usuario:contraseña@host:puerto/basededatos`
*   **Persistencia entre Reinicios:** El contenedor de PostgreSQL utiliza un volumen persistente de Docker (`postgres_data`) montado en `/var/lib/postgresql/data`. Esto asegura que los datos no se pierdan al detener o reiniciar el contenedor, a menos que se ejecute explícitamente una destrucción de volúmenes (`docker compose down -v`).
*   **Aplicación del Esquema:** Las tablas se definen en `schema.sql` y se despliegan automáticamente al inicializar el entorno mediante las tareas definidas en el Makefile.
*   **Acceso a Datos Tipado:** Las operaciones CRUD se centralizan en `queries.sql`. A partir de estas, **sqlc** genera código Go idiomático (`queries.sql.go`), previniendo errores de sintaxis SQL en tiempo de ejecución y aislando la lógica de negocio del lenguaje de base de datos.

---

## 6. Guía de Ejecución y Flujo de Trabajo

### Paso 1: Clonar el Repositorio
Clonar el proyecto desde el repositorio oficial de GitHub:
```bash
git clone https://github.com/tu-usuario/SISTEMA_STOCK.git
cd SISTEMA_STOCK
```

### Paso 2: Ejecución de la Automatización Integral
Parado en la raíz del proyecto, el flujo completo de validación y pruebas se ejecuta con un solo comando gracias al `Makefile`:

```bash
make test
# o el comando principal de integración definido en el Makefile
```

### ¿Qué realiza este comando de automatización?
1.  **Generación de código:** Ejecuta `sqlc generate` para actualizar los modelos y consultas Go a partir de los archivos SQL.
2.  **Limpieza previa:** Detiene y elimina contenedores y volúmenes de Docker anteriores para garantizar un entorno limpio.
3.  **Levantamiento de infraestructura:** Inicia el contenedor de PostgreSQL mediante Docker Compose.
4.  **Healthcheck:** Espera activamente a que la base de datos esté lista para aceptar conexiones.
5.  **Migración de esquemas:** Aplica el archivo `schema.sql` creando las tablas necesarias.
6.  **Ejecución de pruebas:** Corre la suite completa de tests automatizados de la aplicación.
7.  **Limpieza final:** Al terminar los tests (según la configuración del Makefile), desmantela los contenedores y volúmenes temporales para dejar el entorno limpio.