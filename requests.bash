#!/bin/bash
# Prueba cada endpoint de la API de Ctrlstock y verifica el código de respuesta.
# Supone una base recién creada (los IDs arrancan en 1).

BASE="${BASE:-http://localhost:8080}"
PASADOS=0
FALLIDOS=0
RESP=$(mktemp)

# check METODO RUTA CODIGO_ESPERADO [JSON]
check() {
    local metodo="$1" ruta="$2" esperado="$3" datos="${4:-}"
    local codigo

    if [ -n "$datos" ]; then
        codigo=$(curl -s -o "$RESP" -w "%{http_code}" -X "$metodo" "$BASE$ruta" \
            -H "Content-Type: application/json" -d "$datos")
    else
        codigo=$(curl -s -o "$RESP" -w "%{http_code}" -X "$metodo" "$BASE$ruta")
    fi

    if [ "$codigo" = "$esperado" ]; then
        echo "  [OK]    $metodo $ruta -> $codigo"
        PASADOS=$((PASADOS + 1))
    else
        echo "  [FALLA] $metodo $ruta -> $codigo (esperado: $esperado)"
        echo "          respuesta: $(head -c 200 "$RESP")"
        FALLIDOS=$((FALLIDOS + 1))
    fi
}

echo "### CATEGORIA"
check POST   /categoria      201 '{"nombre":"Periféricos"}'
check POST   /categoria      201 '{"nombre":"Teclados","padre_id":1}'
check POST   /categoria      400 '{"nombre":"   "}'                                  # nombre vacío
check GET    /categoria      200
check GET    /categoria/1    200
check GET    /categoria/9999 404
check PUT    /categoria/2    200 '{"nombre":"Teclados mecánicos","padre_id":1}'
check PUT    /categoria/2    400 '{"nombre":"Teclados","padre_id":2}'                # es su propio padre
check PUT    /categoria/9999 404 '{"nombre":"X"}'

echo "### PRODUCTO"
check POST   /producto       201 '{"nombre":"Teclado","descripcion":"Mecánico","precio":45000,"stock":10,"id_categoria":2}'
check POST   /producto       400 '{"nombre":"","precio":10,"stock":1}'               # nombre vacío
check POST   /producto       400 '{"nombre":"Mouse","precio":0,"stock":1}'           # precio
check POST   /producto       400 '{"nombre":"Mouse","precio":100,"stock":-1}'        # stock
check GET    /producto       200
check GET    /producto/1     200
check GET    /producto/9999  404
check PUT    /producto/1     200 '{"nombre":"Teclado RGB","precio":52000,"stock":8,"id_categoria":2}'
check PUT    /producto/9999  404 '{"nombre":"X","precio":1,"stock":1}'

echo "### MOVIMIENTO_STOCK"
check POST   /movimiento_stock      201 '{"producto_id":1,"tipo":"ingreso","cantidad":5,"motivo":"Compra a proveedor"}'
check POST   /movimiento_stock      400 '{"producto_id":1,"tipo":"OTRO","cantidad":5,"motivo":"x"}'     # tipo
check POST   /movimiento_stock      400 '{"producto_id":1,"tipo":"EGRESO","cantidad":0,"motivo":"x"}'   # cantidad
check GET    /movimiento_stock      200
check GET    /movimiento_stock/1    200
check GET    /movimiento_stock/9999 404

echo "### DELETE (primero los que dependen de otros)"
check DELETE /movimiento_stock/1 204
check DELETE /producto/1         204
check DELETE /categoria/2        204
check DELETE /categoria/1        204

rm -f "$RESP"
echo
echo "Resultado: $PASADOS pasaron, $FALLIDOS fallaron"
[ "$FALLIDOS" -eq 0 ]