#!/usr/bin/env bash
# Pruebas de la aplicación renderizada en el servidor (SSR).
# Requiere el servidor corriendo y la base recién creada.

URL="${APP_URL:-http://localhost:8080}"
ok=0
falla=0

verificar() {
    local descripcion="$1" esperado="$2" obtenido="$3"
    if [ "$esperado" == "$obtenido" ]; then
        echo "[OK]    $descripcion"
        ok=$((ok + 1))
    else
        echo "[FALLA] $descripcion (esperaba $esperado, obtuve $obtenido)"
        falla=$((falla + 1))
    fi
}

# 1. La página principal responde con HTML
codigo=$(curl -s -o /dev/null -w "%{http_code}" "$URL/")
verificar "GET / responde 200" "200" "$codigo"

tipo=$(curl -s -o /dev/null -w "%{content_type}" "$URL/")
verificar "GET / devuelve HTML" "text/html; charset=utf-8" "$tipo"

# 2. Con la base vacía se muestra el mensaje correspondiente
if curl -s "$URL/" | grep -q "Todavía no hay productos"; then vacio="si"; else vacio="no"; fi
verificar "GET / muestra el mensaje de lista vacía" "si" "$vacio"

# 3. Enviar el formulario redirige a / (patrón PRG)
codigo=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/productos" \
    --data-urlencode "nombre=Producto de prueba" \
    --data-urlencode "descripcion=Creado desde requests.bash" \
    --data-urlencode "precio=1500.50" \
    --data-urlencode "stock=10" \
    --data-urlencode "id_categoria=")
verificar "POST /productos responde 303 (redirect)" "303" "$codigo"

destino=$(curl -s -o /dev/null -w "%{redirect_url}" -X POST "$URL/productos" \
    --data-urlencode "nombre=Otro producto" \
    --data-urlencode "precio=200" \
    --data-urlencode "stock=1")
verificar "POST /productos redirige a /" "$URL/" "$destino"

# 4. Después del redirect, la lista muestra el producto nuevo
if curl -s "$URL/" | grep -q "Producto de prueba"; then aparece="si"; else aparece="no"; fi
verificar "GET / muestra el producto creado" "si" "$aparece"

# 5. Datos inválidos se rechazan
codigo=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/productos" \
    --data-urlencode "nombre=Sin precio" \
    --data-urlencode "precio=0" \
    --data-urlencode "stock=1")
verificar "POST /productos con precio 0 responde 400" "400" "$codigo"

codigo=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/productos" \
    --data-urlencode "nombre=" \
    --data-urlencode "precio=100" \
    --data-urlencode "stock=1")
verificar "POST /productos sin nombre responde 400" "400" "$codigo"

echo
echo "Resultado: $ok pasaron, $falla fallaron"
[ "$falla" -eq 0 ]