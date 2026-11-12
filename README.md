# CtrlStock: gestión de stock con renderizado en el servidor

**CtrlStock** es una aplicación web para gestionar el stock de un depósito. Permite administrar productos, organizarlos en categorías (con jerarquía), registrar los movimientos de inventario y ver qué productos se venden más.

Es una aplicación **renderizada en el servidor (SSR)**: Go genera el HTML con **templ** y lo envía ya armado al navegador. No hay una API JSON ni JavaScript del lado del cliente. Los formularios se envían con `POST` y el servidor responde con una redirección, siguiendo el patrón **Post/Redirect/Get (PRG)**.

---

## 1. Tecnologías utilizadas

*   **Go:** lenguaje de la aplicación y del servidor HTTP (`net/http`).
*   **templ:** genera las vistas HTML como componentes Go, verificados en compilación.
*   **Pico.css:** hoja de estilos (cargada por CDN), más estilos propios definidos en el layout.
*   **PostgreSQL:** motor de base de datos relacional.
*   **Docker y Docker Compose:** levantan PostgreSQL en un contenedor, sin instalaciones locales.
*   **sqlc:** genera código Go tipado a partir de las consultas SQL. Ese código no se versiona: se crea al ejecutar `make run`.
*   **lib/pq:** driver de PostgreSQL para Go.
*   **GNU Make:** automatiza la generación de código, la base de datos y la ejecución.

---

## 2. Requisitos

1.  **Go** `>= 1.22.2` (las rutas usan patrones con método, como `GET /{$}`)
2.  **Docker** y **Docker Compose**
3.  **Make**
4.  **sqlc** y **templ**, instalados con `go install` y disponibles en el `PATH`:

    ```bash
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    go install github.com/a-h/templ/cmd/templ@latest
    export PATH="$PATH:$HOME/go/bin"   # si todavía no está en el PATH
    ```
5.  **curl** (solo si se quiere correr `make test`)
6.  Conexión a internet (para descargar Pico.css y las dependencias de Go).

---

## 3. Cómo ejecutar el proyecto

### Paso 1: Clonar el repositorio

```bash
git clone https://github.com/fioreferreri-lab/Ctrlstock.git
cd Ctrlstock
git checkout tp5
```

### Paso 2: Levantar la aplicación

```bash
make run
```

Este comando, en orden:

1.  levanta PostgreSQL en Docker y espera a que esté lista,
2.  crea las tablas a partir de `db/schema/schema.sql`,
3.  genera el código Go: las consultas (**sqlc**) y las vistas (**templ**),
4.  inicia el servidor.

La aplicación queda disponible en **http://localhost:8080**. Para detenerla: `Ctrl + C`.

> **Importante:** cada `make run` vuelve a crear las tablas, así que **los datos no se conservan entre ejecuciones**: se arranca siempre con la base vacía.

La conexión a la base se lee de la variable de entorno `DATABASE_URL`, que el `makefile` define automáticamente (debe coincidir con `docker-compose.yaml`):

```text
postgres://user:password@localhost:5432/stockapp?sslmode=disable
```

### Otros comandos del makefile

| Comando | Descripción |
|---|---|
| `make run` | Levanta la base, genera el código y arranca el servidor |
| `make generate` | Ejecuta `sqlc generate` y `templ generate` |
| `make build` | Genera el código y compila el binario en `tmp/` |
| `make test` | Ciclo completo de pruebas sobre una base limpia (ver sección 8) |
| `make db-reset` | Borra el contenedor y el volumen de la base (datos incluidos) |
| `make clean` | Borra la base y los archivos compilados |

---

## 4. Páginas de la aplicación

Todas las páginas comparten el layout con el menú lateral.

| Ruta | Página | Qué permite hacer |
|---|---|---|
| `/` | — | Redirige a `/producto` |
| `/deposito` | **Depósito** | Ver todo lo que hay en stock organizado por categoría (árbol desplegable), con las unidades en stock, la inversión total y la cantidad de categorías principales |
| `/producto` | **Productos** | Ver el inventario y cargar un producto nuevo (nombre, descripción, precio, stock inicial y categoría) |
| `/categoria` | **Categorías** | Ver las categorías y crear una nueva, con una categoría padre opcional |
| `/movimiento_stock` | **Movimientos** | Ver el historial y registrar un ingreso o un egreso de un producto |
| `/ventas?dias=7\|30\|90` | **Más vendidos** | Ranking de productos por unidades vendidas en los últimos 7, 30 (por defecto) o 90 días |

Las rutas `/producto/{id}` y `/categoria/{id}` permiten editar (`POST`/`PUT`) y eliminar (`DELETE`) un registro.

Cuando un formulario tiene un error de validación, la página se vuelve a mostrar con el mensaje y con los datos que el usuario ya había escrito (respuesta `400`). Cuando todo sale bien, el servidor responde con una redirección `303` y el navegador pide la página de nuevo con la lista actualizada.

### Flujo de una petición

```text
Navegador → ruta → handler → reglas de negocio (logic) → consultas sqlc → PostgreSQL
                      │
                      └→ componente templ (views) → HTML completo → Navegador
```

---

## 5. Reglas de negocio

Las validan los archivos de `logic/`, no las vistas, para que valgan sin importar desde dónde llegue el pedido.

**Productos**
*   Todo producto debe tener una **categoría obligatoria** y existente: no hay productos sin categoría.
*   **No puede haber nombres repetidos.** Para sumar unidades de un producto que ya existe se registra un ingreso en Movimientos.
*   El nombre no puede estar vacío, el precio debe ser mayor a cero y el stock no puede ser negativo.
*   No se puede eliminar un producto que tiene movimientos registrados.

**Categorías**
*   El nombre es obligatorio, de hasta 100 caracteres, y **no puede repetirse**.
*   Una categoría no puede ser su propia categoría padre.
*   No se puede eliminar una categoría que tiene productos.

**Stock y movimientos**
*   El stock se define una sola vez, al crear el producto (stock inicial). **Después solo cambia con movimientos:** un ingreso suma unidades y un egreso las resta. El stock no se edita desde el formulario de producto.
*   Un egreso no puede dejar el stock en negativo (`stock insuficiente`).
*   Registrar un movimiento y actualizar el stock ocurre en **una sola transacción**: se guardan las dos cosas o ninguna.
*   **Los movimientos no se modifican ni se eliminan:** son el historial de stock. Un error de carga se corrige registrando el movimiento contrario. Los pedidos para borrar un movimiento responden `405`.
*   El motivo se elige de una lista fija, según el tipo:
    *   **Ingreso:** Compra a proveedor, Devolución de cliente, Ajuste de inventario.
    *   **Egreso:** Venta, Rotura o pérdida, Ajuste de inventario.

**Más vendidos**
*   Solo cuentan los egresos con motivo **Venta**. Las roturas y los ajustes no suman al ranking.

---

## 6. Estructura del proyecto

```text
Ctrlstock/
│
├── main.go                    # Conexión a la base, rutas y servidor HTTP
│
├── handlers/                  # Capa HTTP: leen el formulario, llaman a logic y renderizan
│   ├── products.go
│   ├── categories.go
│   ├── stock_movements.go
│   ├── deposito.go
│   ├── ventas.go
│   └── utils.go               # Lectura y chequeo básico de los campos del formulario
│
├── logic/                     # Reglas de negocio (validaciones, transacción de movimientos, ranking)
│   ├── product.go
│   ├── category.go
│   └── stock_movement.go
│
├── views/                     # Componentes templ (*.templ) y el código Go que generan (*_templ.go)
│   ├── layout.templ           # Estructura base: HTML, head con Pico.css y menú lateral
│   └── ...                    # Formulario, lista y página de cada sección
│
├── db/                        # Capa de datos
│   ├── schema/schema.sql      # Definición de las tablas
│   ├── queries/               # Consultas SQL para sqlc (un archivo por tabla)
│   ├── *.go generados         # Código de sqlc (no se sube al repo)
│   └── crud_test.go           # Tests de las consultas contra PostgreSQL
│
├── docker-compose.yaml        # Contenedor de PostgreSQL
├── sqlc.yaml                  # Configuración de sqlc
├── makefile                   # Automatización de tareas
└── requests.bash              # Pruebas con curl de las páginas, formularios y reglas
```

---

## 7. Modelo de datos

*   **Categoría (`categoria`):** nombre único y categoría padre opcional (`padre_id`), que arma la jerarquía.
*   **Producto (`producto`):** nombre único, descripción opcional, precio, stock actual y categoría (`id_categoria`, obligatoria).
*   **Movimiento de stock (`movimiento_stock`):** producto, tipo (`INGRESO` o `EGRESO`), cantidad, motivo y fecha.

---

## 8. Pruebas (opcional)

```bash
make test
```

Parte de una base limpia y realiza, en orden:

1.  **Tareas previas:** genera el código, compila, recrea la base y levanta la aplicación.
2.  **Pruebas:** corre los tests de Go de la capa de datos y luego `requests.bash`, que con `curl`:
    *   verifica que cada página responde con HTML,
    *   envía los formularios y comprueba la redirección `303` y que el dato aparece en la lista,
    *   comprueba las reglas de negocio (categoría obligatoria, nombres repetidos, stock insuficiente, motivos válidos, movimientos que no se borran, etc.).
3.  **Tareas posteriores:** apaga la aplicación y borra contenedores y volúmenes, aunque alguna prueba falle.

Cada prueba se muestra como `[OK]` o `[FALLA]`, y al final se imprime un resumen (`Resultado: N pasaron, 0 fallaron`).

También se puede correr `requests.bash` por separado, con la aplicación ya levantada (`make run`):

```bash
bash ./requests.bash
```