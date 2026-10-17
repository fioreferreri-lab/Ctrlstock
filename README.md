# CtrlStock: API REST de gestión de stock

**CtrlStock** es una API REST para la gestión de stock e inventario. Permite administrar productos, una jerarquía de categorías y el registro histórico de los movimientos de inventario.

La API expone la lógica de negocio a través de endpoints HTTP, conectándola con la capa de datos generada por **sqlc** sobre **PostgreSQL**.

---

## 1. Tecnologías utilizadas

*   **Go:** lenguaje principal de la aplicación y del servidor HTTP (`net/http`).
*   **PostgreSQL:** motor de base de datos relacional.
*   **Docker y Docker Compose:** levantan PostgreSQL en un contenedor, sin instalaciones locales y de forma reproducible en cualquier entorno.
*   **sqlc:** genera código Go tipado a partir de las consultas SQL. El código generado no se versiona: se crea automáticamente al ejecutar `make run`, `make build` o `make test`.
*   **lib/pq:** driver de PostgreSQL para Go.
*   **GNU Make:** automatiza la generación de código, la base de datos, la ejecución y las pruebas.
*   **curl:** usado por `requests.bash` para probar cada endpoint.

---

## 2. Requisitos

1.  **Go** `>= 1.22.2`
2.  **Docker** y **Docker Compose**
3.  **Make**
4.  **sqlc** (debe estar en el `PATH`; por ejemplo, si se instaló con `go install`, agregar `$HOME/go/bin` al `PATH`)
5.  **curl**

---

## 3. Estructura del proyecto

```text
Ctrlstock/
│
├── main.go                    # Conexión a la base, instancia de sqlc, rutas y servidor HTTP
│
├── handlers/                  # Capa HTTP: reciben el pedido, validan y responden JSON
│   ├── categories.go
│   ├── products.go
│   ├── stock_movements.go
│   └── utils.go
│
├── logic/                     # Reglas de negocio (validaciones)
│   ├── category.go
│   ├── product.go
│   └── stock_movement.go
│
├── db/                        # Capa de datos
│   ├── scheme/schema.sql      # Definición de las tablas
│   ├── queries/               # Consultas SQL usadas por sqlc (un archivo por tabla)
│   │   ├── categoria.sql
│   │   ├── producto.sql
│   │   └── movimiento_stock.sql
│   ├── *.go generados         # Código generado por sqlc (no se sube al repo)
│   └── crud_test.go           # Tests de las consultas contra PostgreSQL
│
├── docker-compose.yaml        # Contenedor de PostgreSQL
├── sqlc.yaml                  # Configuración de sqlc
├── makefile                   # Automatización de tareas
└── requests.bash              # Pruebas de todos los endpoints con curl
```

### Flujo de una petición

```text
Cliente (curl) → endpoint → handler → reglas de negocio (logic) → queries de sqlc → PostgreSQL
                                                                                      ↓
Cliente ← respuesta JSON ← handler ←───────────────────────────────────────────────  datos
```

---

## 4. Modelo de datos

### A. Categoría (`categoria`)
Organiza los productos en una jerarquía: una categoría puede tener una categoría padre (`padre_id`, opcional).

### B. Producto (`producto`)
Cada ítem del inventario: nombre, descripción, precio, stock actual y categoría (`id_categoria`).

### C. Movimiento de stock (`movimiento_stock`)
Registro histórico de cada movimiento de inventario: producto, tipo (`INGRESO` o `EGRESO`), cantidad, motivo y fecha.

---

## 5. Cómo ejecutar el proyecto

### Paso 1: Clonar el repositorio

```bash
git clone https://github.com/fioreferreri-lab/Ctrlstock.git
cd Ctrlstock
git checkout tp3-api-rest
```

### Paso 2: Levantar el servidor

```bash
make run
```

Este comando:
1.  levanta PostgreSQL en Docker y espera a que esté lista,
2.  crea las tablas a partir de `schema.sql`,
3.  genera el código con sqlc,
4.  inicia el servidor.

La API queda disponible en **http://localhost:8080**. Para detenerla: `Ctrl + C`.

La aplicación lee la conexión a la base desde la variable de entorno `DATABASE_URL`, que el `makefile` define automáticamente:

```text
postgres://user:password@localhost:5432/stockapp?sslmode=disable
```

### Paso 3: Probar la API

Con un solo comando, desde una base limpia:

```bash
make test
```

Este comando realiza, en orden:

1.  **Tareas previas:** genera el código con sqlc, compila la aplicación, borra contenedores y volúmenes anteriores, levanta PostgreSQL, espera a que esté lista, crea las tablas y levanta la aplicación.
2.  **Pruebas:** corre los tests de Go de la capa de datos y luego ejecuta `requests.bash`, que prueba cada endpoint y verifica el código de respuesta.
3.  **Tareas posteriores:** apaga la aplicación y borra contenedores y volúmenes (aunque alguna prueba falle).

Cada prueba se muestra como `[OK]` o `[FALLA]`, y al final se imprime un resumen:

```text
Resultado: 28 pasaron, 0 fallaron
```

También se puede correr `requests.bash` por separado, con el servidor ya levantado (`make run`) y la base recién creada:

```bash
./requests.bash
```

### Otros comandos del makefile

| Comando | Descripción |
|---|---|
| `make run` | Levanta la base y el servidor |
| `make test` | Ciclo completo de pruebas sobre una base limpia |
| `make build` | Genera el código y compila el binario en `tmp/` |
| `make generate` | Ejecuta `sqlc generate` |
| `make db-reset` | Borra contenedores y volúmenes (datos incluidos) |
| `make clean` | Borra la base y los archivos compilados |

---

## 6. Endpoints de la API

Todas las entidades exponen las mismas operaciones CRUD:

| Método | Ruta | Descripción | Respuesta |
|---|---|---|---|
| `POST` | `/<entidad>` | Crea un registro a partir de un JSON validado | `201` con el objeto creado · `400` si los datos son inválidos |
| `GET` | `/<entidad>` | Lista todos los registros | `200` con un array JSON |
| `GET` | `/<entidad>/{id}` | Obtiene un registro por ID | `200` · `404` si no existe |
| `PUT` | `/<entidad>/{id}` | Actualiza un registro con datos validados | `200` con el objeto actualizado · `400` · `404` |
| `DELETE` | `/<entidad>/{id}` | Elimina un registro | `204 No Content` |

Entidades: `categoria`, `producto` y `movimiento_stock`.

### Categoría

```bash
curl -X POST http://localhost:8080/categoria \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Teclados", "padre_id": 1}'
```

Validaciones: el nombre no puede estar vacío y una categoría no puede ser su propia categoría padre.

### Producto

```bash
curl -X POST http://localhost:8080/producto \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Teclado", "descripcion": "Mecánico", "precio": 45000, "stock": 10, "id_categoria": 2}'
```

Validaciones: el nombre no puede estar vacío, el precio debe ser mayor a cero y el stock no puede ser negativo.

### Movimiento de stock

```bash
curl -X POST http://localhost:8080/movimiento_stock \
  -H "Content-Type: application/json" \
  -d '{"producto_id": 1, "tipo": "INGRESO", "cantidad": 5, "motivo": "Compra a proveedor"}'
```

Validaciones: el tipo debe ser `INGRESO` o `EGRESO` y la cantidad debe ser mayor a cero.

### Errores

Los errores de validación responden `400 Bad Request` y los registros inexistentes `404 Not Found`, con un mensaje descriptivo en el cuerpo de la respuesta (por ejemplo, `el precio debe ser mayor a cero`).

El archivo `requests.bash` contiene ejemplos de todas las peticiones probadas, incluyendo los casos de error.