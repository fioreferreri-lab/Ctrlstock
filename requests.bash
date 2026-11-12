#!/usr/bin/env bash
# Pruebas de la aplicación renderizada en el servidor (SSR con templ).
# Verifica las páginas, los formularios (patrón Post/Redirect/Get) y las reglas de negocio.
#
# Requiere la aplicación corriendo (make run, o make test que la levanta solo).
# Usa nombres únicos en cada corrida; el chequeo del ranking de "Más vendidos"
# asume que el producto de prueba entra en el ranking (base limpia, como en make test).

URL="${APP_URL:-http://localhost:8080}"
SUF="$(date +%s)"
CAT="Categoria test $SUF"
PROD="Producto test $SUF"

ok=0
falla=0
CODIGO=""
CUERPO=""

# ---------------------------------------------------------------
# Funciones de apoyo
# ---------------------------------------------------------------

# pedir <argumentos de curl>: guarda el código HTTP en $CODIGO y el cuerpo en $CUERPO.
# No sigue redirecciones, así que un POST exitoso se ve como 303.
pedir() {
	local salida
	salida=$(curl -s -w $'\n%{http_code}' "$@")
	CODIGO="${salida##*$'\n'}"
	CUERPO="${salida%$'\n'*}"
}

registrar() {
	if [ "$1" == "ok" ]; then
		echo "[OK]    $2"
		ok=$((ok + 1))
	else
		echo "[FALLA] $2"
		falla=$((falla + 1))
	fi
}

# verificar <descripción> <esperado> <obtenido>
verificar() {
	if [ "$2" == "$3" ]; then
		registrar ok "$1"
	else
		registrar falla "$1 (esperaba $2, obtuve $3)"
	fi
}

# verificar_contiene <descripción> <texto>: busca el texto en el último cuerpo recibido.
verificar_contiene() {
	if grep -qiF -- "$2" <<<"$CUERPO"; then
		registrar ok "$1"
	else
		registrar falla "$1 (no aparece: $2)"
	fi
}

# verificar_4xx <descripción>: el último pedido tiene que haber sido rechazado (400 a 499).
verificar_4xx() {
	case "$CODIGO" in
	4??) registrar ok "$1" ;;
	*) registrar falla "$1 (esperaba 4xx, obtuve $CODIGO)" ;;
	esac
}

# id_opcion <html> <texto>: devuelve el value del <option> cuyo texto es <texto>.
id_opcion() {
	tr '\n' ' ' <<<"$1" |
		grep -oE "<option value=\"[0-9]+\"[^>]*>[[:space:]]*$2[[:space:]]*</option>" |
		head -1 | grep -oE 'value="[0-9]+"' | grep -oE '[0-9]+'
}

# movimiento <tipo> <cantidad> <motivo> [producto_id]: envía el formulario de movimientos.
movimiento() {
	pedir -X POST "$URL/movimiento_stock" \
		--data-urlencode "producto_id=${4:-$ID_PROD}" \
		--data-urlencode "tipo=$1" \
		--data-urlencode "cantidad=$2" \
		--data-urlencode "motivo=$3"
}

# ---------------------------------------------------------------
# 1. Páginas
# ---------------------------------------------------------------
echo "== Páginas =="

destino=$(curl -s -o /dev/null -w "%{redirect_url}" "$URL/")
codigo=$(curl -s -o /dev/null -w "%{http_code}" "$URL/")
verificar "GET / responde 303 (redirección)" "303" "$codigo"
verificar "GET / redirige a /producto" "$URL/producto" "$destino"

for ruta in /producto /categoria /movimiento_stock /deposito /ventas; do
	tipo=$(curl -s -o /dev/null -w "%{content_type}" "$URL$ruta")
	pedir "$URL$ruta"
	verificar "GET $ruta responde 200" "200" "$CODIGO"
	case "$tipo" in
	text/html*) registrar ok "GET $ruta devuelve HTML" ;;
	*) registrar falla "GET $ruta devuelve HTML (obtuve $tipo)" ;;
	esac
done

pedir "$URL/ventas?dias=7"
verificar "GET /ventas?dias=7 responde 200" "200" "$CODIGO"
pedir "$URL/ventas?dias=90"
verificar "GET /ventas?dias=90 responde 200" "200" "$CODIGO"
pedir "$URL/ventas?dias=abc"
verificar "GET /ventas con un período inválido usa el valor por defecto (200)" "200" "$CODIGO"

# ---------------------------------------------------------------
# 2. Categorías
# ---------------------------------------------------------------
echo
echo "== Categorías =="

pedir -X POST "$URL/categoria" --data-urlencode "nombre=$CAT"
verificar "POST /categoria responde 303 (PRG)" "303" "$CODIGO"

pedir "$URL/categoria"
verificar_contiene "GET /categoria muestra la categoría creada" "$CAT"

pedir -X POST "$URL/categoria" --data-urlencode "nombre=$CAT"
verificar_4xx "POST /categoria con nombre repetido se rechaza"

pedir -X POST "$URL/categoria" --data-urlencode "nombre="
verificar_4xx "POST /categoria sin nombre se rechaza"

# ---------------------------------------------------------------
# 3. Productos
# ---------------------------------------------------------------
echo
echo "== Productos =="

pedir "$URL/producto"
ID_CAT="$(id_opcion "$CUERPO" "$CAT")"
if [ -n "$ID_CAT" ]; then
	registrar ok "la categoría creada aparece en el desplegable del producto (id $ID_CAT)"
else
	registrar falla "la categoría creada aparece en el desplegable del producto"
	ID_CAT=0
fi

# Regla: categoría obligatoria y existente
pedir -X POST "$URL/producto" --data-urlencode "nombre=Sin categoria $SUF" \
	--data-urlencode "precio=100" --data-urlencode "stock=1" --data-urlencode "id_categoria="
verificar "POST /producto sin categoría responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que falta la categoría" "debe tener una categoría"

pedir -X POST "$URL/producto" --data-urlencode "nombre=Categoria falsa $SUF" \
	--data-urlencode "precio=100" --data-urlencode "stock=1" --data-urlencode "id_categoria=999999"
verificar "POST /producto con una categoría inexistente responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que la categoría no existe" "no existe"

# Reglas: datos inválidos
pedir -X POST "$URL/producto" --data-urlencode "nombre=Precio cero $SUF" \
	--data-urlencode "precio=0" --data-urlencode "stock=1" --data-urlencode "id_categoria=$ID_CAT"
verificar "POST /producto con precio 0 responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que el precio debe ser mayor a cero" "mayor a cero"

pedir -X POST "$URL/producto" --data-urlencode "nombre=" \
	--data-urlencode "precio=100" --data-urlencode "stock=1" --data-urlencode "id_categoria=$ID_CAT"
verificar "POST /producto sin nombre responde 400" "400" "$CODIGO"

pedir -X POST "$URL/producto" --data-urlencode "nombre=Stock negativo $SUF" \
	--data-urlencode "precio=100" --data-urlencode "stock=-1" --data-urlencode "id_categoria=$ID_CAT"
verificar "POST /producto con stock negativo responde 400" "400" "$CODIGO"

# Alta correcta (stock inicial 10)
pedir -X POST "$URL/producto" --data-urlencode "nombre=$PROD" \
	--data-urlencode "descripcion=Creado desde requests.bash" \
	--data-urlencode "precio=1500.50" --data-urlencode "stock=10" \
	--data-urlencode "id_categoria=$ID_CAT"
verificar "POST /producto responde 303 (PRG)" "303" "$CODIGO"

pedir "$URL/producto"
verificar_contiene "GET /producto muestra el producto creado" "$PROD"
verificar_contiene "GET /producto muestra el precio del producto" "1500.50"

# Regla: nombres repetidos
pedir -X POST "$URL/producto" --data-urlencode "nombre=$PROD" \
	--data-urlencode "precio=100" --data-urlencode "stock=1" --data-urlencode "id_categoria=$ID_CAT"
verificar "POST /producto con nombre repetido responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que ya existe" "ya existe un producto"

# ---------------------------------------------------------------
# 4. Movimientos y stock (el producto arranca con 10 unidades)
# ---------------------------------------------------------------
echo
echo "== Movimientos =="

pedir "$URL/movimiento_stock"
ID_PROD="$(id_opcion "$CUERPO" "$PROD")"
if [ -n "$ID_PROD" ]; then
	registrar ok "el producto creado aparece en el desplegable de movimientos (id $ID_PROD)"
else
	registrar falla "el producto creado aparece en el desplegable de movimientos"
	ID_PROD=0
fi

# Validaciones del formulario
movimiento INGRESO 1 "Venta"
verificar "un ingreso con motivo de egreso responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que el motivo no corresponde" "no corresponde a un ingreso"

movimiento OTRO 1 "Venta"
verificar "un tipo de movimiento inválido responde 400" "400" "$CODIGO"

movimiento EGRESO 0 "Venta"
verificar "un movimiento con cantidad 0 responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que la cantidad debe ser mayor a cero" "mayor a cero"

movimiento EGRESO 1 "Venta" 999999
verificar "un movimiento de un producto inexistente responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje explica que el producto no existe" "no existe"

# El stock cambia con los movimientos: 10 -> 6 -> 0
movimiento EGRESO 4 "Venta"
verificar "egreso de 4 unidades (stock 10 -> 6) responde 303" "303" "$CODIGO"

movimiento EGRESO 7 "Venta"
verificar "egreso de 7 unidades con 6 en stock responde 400" "400" "$CODIGO"
verificar_contiene "...y el mensaje indica stock insuficiente" "stock insuficiente"

movimiento EGRESO 6 "Venta"
verificar "egreso de 6 unidades (stock 6 -> 0) responde 303" "303" "$CODIGO"

movimiento EGRESO 1 "Venta"
verificar "egreso de 1 unidad con 0 en stock responde 400" "400" "$CODIGO"

# Ingreso y egreso por otros motivos: 0 -> 5 -> 0
movimiento INGRESO 5 "Compra a proveedor"
verificar "ingreso de 5 unidades (stock 0 -> 5) responde 303" "303" "$CODIGO"

movimiento EGRESO 5 "Rotura o pérdida"
verificar "egreso por rotura de 5 unidades (stock 5 -> 0) responde 303" "303" "$CODIGO"

# Regla: el stock no se edita desde el producto, solo con movimientos
pedir -X POST "$URL/producto/$ID_PROD" --data-urlencode "nombre=$PROD" \
	--data-urlencode "precio=1600" --data-urlencode "stock=999" \
	--data-urlencode "id_categoria=$ID_CAT"
verificar "editar el producto responde 303" "303" "$CODIGO"

movimiento EGRESO 1 "Venta"
verificar "editar el producto no cambia el stock (sigue en 0)" "400" "$CODIGO"

# Historial
pedir "$URL/movimiento_stock"
verificar_contiene "el historial muestra el producto de los movimientos" "$PROD"
verificar_contiene "el historial muestra el motivo Venta" "Venta"
verificar_contiene "el historial muestra el motivo Rotura o pérdida" "Rotura o pérdida"

# Regla: los movimientos no se borran
pedir -X DELETE "$URL/movimiento_stock/1"
verificar "DELETE /movimiento_stock/1 responde 405" "405" "$CODIGO"
pedir -X POST "$URL/movimiento_stock/1"
verificar "POST /movimiento_stock/1 responde 405" "405" "$CODIGO"

# ---------------------------------------------------------------
# 5. Borrados que protegen la integridad de los datos
# ---------------------------------------------------------------
echo
echo "== Borrados =="

pedir -X DELETE "$URL/producto/$ID_PROD"
verificar "borrar un producto con movimientos responde 409" "409" "$CODIGO"

pedir -X DELETE "$URL/categoria/$ID_CAT"
verificar "borrar una categoría con productos responde 409" "409" "$CODIGO"

# ---------------------------------------------------------------
# 6. Depósito y Más vendidos
# ---------------------------------------------------------------
echo
echo "== Depósito y Más vendidos =="

pedir "$URL/deposito"
verificar_contiene "el Depósito muestra la categoría" "$CAT"
verificar_contiene "el Depósito muestra el producto" "$PROD"

pedir "$URL/ventas?dias=30"
verificar_contiene "Más vendidos muestra el producto vendido" "$PROD"

# ---------------------------------------------------------------
echo
echo "Resultado: $ok pasaron, $falla fallaron"
[ "$falla" -eq 0 ]