#!/bin/bash
BASE=http://localhost:8080

echo "### CATEGORIA"
curl -i -X POST $BASE/categoria -H "Content-Type: application/json" \
  -d '{"nombre":"Periféricos"}'                                      # 201
curl -i -X POST $BASE/categoria -H "Content-Type: application/json" \
  -d '{"nombre":"Teclados","padre_id":1}'                            # 201
curl -i -X POST $BASE/categoria -H "Content-Type: application/json" \
  -d '{"nombre":"   "}'                                              # 400 nombre vacío
curl -i $BASE/categoria                                              # 200
curl -i $BASE/categoria/1                                            # 200
curl -i $BASE/categoria/9999                                         # 404
curl -i -X PUT $BASE/categoria/2 -H "Content-Type: application/json" \
  -d '{"nombre":"Teclados mecánicos","padre_id":1}'                  # 200
curl -i -X PUT $BASE/categoria/2 -H "Content-Type: application/json" \
  -d '{"nombre":"Teclados","padre_id":2}'                            # 400 es su propio padre
curl -i -X PUT $BASE/categoria/9999 -H "Content-Type: application/json" \
  -d '{"nombre":"X"}'                                                # 404

echo "### PRODUCTO"
curl -i -X POST $BASE/producto -H "Content-Type: application/json" \
  -d '{"nombre":"Teclado","descripcion":"Mecánico","precio":45000,"stock":10,"id_categoria":2}'  # 201
curl -i -X POST $BASE/producto -H "Content-Type: application/json" \
  -d '{"nombre":"","precio":10,"stock":1}'                           # 400 nombre vacío
curl -i -X POST $BASE/producto -H "Content-Type: application/json" \
  -d '{"nombre":"Mouse","precio":0,"stock":1}'                       # 400 precio
curl -i -X POST $BASE/producto -H "Content-Type: application/json" \
  -d '{"nombre":"Mouse","precio":100,"stock":-1}'                    # 400 stock
curl -i $BASE/producto                                               # 200
curl -i $BASE/producto/1                                             # 200
curl -i $BASE/producto/9999                                          # 404
curl -i -X PUT $BASE/producto/1 -H "Content-Type: application/json" \
  -d '{"nombre":"Teclado RGB","precio":52000,"stock":8,"id_categoria":2}'  # 200
curl -i -X PUT $BASE/producto/9999 -H "Content-Type: application/json" \
  -d '{"nombre":"X","precio":1,"stock":1}'                           # 404

echo "### MOVIMIENTO_STOCK"
curl -i -X POST $BASE/movimiento_stock -H "Content-Type: application/json" \
  -d '{"producto_id":1,"tipo":"ingreso","cantidad":5,"motivo":"Compra a proveedor"}'  # 201
curl -i -X POST $BASE/movimiento_stock -H "Content-Type: application/json" \
  -d '{"producto_id":1,"tipo":"OTRO","cantidad":5,"motivo":"x"}'     # 400 tipo
curl -i -X POST $BASE/movimiento_stock -H "Content-Type: application/json" \
  -d '{"producto_id":1,"tipo":"EGRESO","cantidad":0,"motivo":"x"}'   # 400 cantidad
curl -i $BASE/movimiento_stock                                       # 200
curl -i $BASE/movimiento_stock/1                                     # 200
curl -i $BASE/movimiento_stock/9999                                  # 404

echo "### DELETE (primero los que dependen de otros)"
curl -i -X DELETE $BASE/movimiento_stock/1                           # 204
curl -i -X DELETE $BASE/producto/1                                   # 204
curl -i -X DELETE $BASE/categoria/2                                  # 204
curl -i -X DELETE $BASE/categoria/1                                  # 204
