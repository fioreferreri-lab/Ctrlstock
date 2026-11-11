package logic

import (
	"strconv"

	"ctrlstock/db"
)

// NodoCategoria es una rama del árbol del depósito: una categoría con sus
// subcategorías, sus productos y los totales acumulados (incluyendo las subcategorías).
type NodoCategoria struct {
	Nombre    string
	Hijas     []*NodoCategoria
	Productos []db.Producto
	Unidades  int64   // unidades de esta categoría + todas sus subcategorías
	Inversion float64 // stock × precio, sumado en todo el subárbol
}

// Deposito es el árbol completo con los totales generales.
type Deposito struct {
	Raices    []*NodoCategoria // categorías principales (y "Sin categoría", si hiciera falta)
	Unidades  int64
	Inversion float64
}

// PrecioDe devuelve el precio del producto como número.
// sqlc lo trae como texto porque en la base es DECIMAL.
func PrecioDe(p db.Producto) float64 {
	precio, err := strconv.ParseFloat(p.Precio, 64)
	if err != nil {
		return 0
	}
	return precio
}

// InversionDe calcula cuánto hay invertido en un producto: stock × precio.
func InversionDe(p db.Producto) float64 {
	return PrecioDe(p) * float64(p.Stock)
}

// ArmarDeposito organiza categorías y productos en un árbol y calcula los totales.
func ArmarDeposito(categorias []db.Categoria, productos []db.Producto) Deposito {
	// 1. Un nodo por categoría, indexado por ID
	nodos := make(map[int32]*NodoCategoria, len(categorias))
	for _, c := range categorias {
		nodos[c.IDCategoria] = &NodoCategoria{Nombre: c.Nombre}
	}

	// 2. Colgar cada categoría de su padre (o dejarla como principal)
	var raices []*NodoCategoria
	for _, c := range categorias {
		nodo := nodos[c.IDCategoria]
		padre, existePadre := nodos[c.PadreID.Int32]
		if c.PadreID.Valid && existePadre && c.PadreID.Int32 != c.IDCategoria {
			padre.Hijas = append(padre.Hijas, nodo)
		} else {
			raices = append(raices, nodo)
		}
	}

	// 3. Ubicar cada producto en su categoría
	sinCategoria := &NodoCategoria{Nombre: "Sin categoría"}
	for _, p := range productos {
		if nodo, ok := nodos[p.IDCategoria.Int32]; p.IDCategoria.Valid && ok {
			nodo.Productos = append(nodo.Productos, p)
		} else {
			sinCategoria.Productos = append(sinCategoria.Productos, p)
		}
	}

	// 4. Calcular los totales de abajo hacia arriba
	visitados := make(map[*NodoCategoria]bool)
	for _, r := range raices {
		calcularTotales(r, visitados)
	}

	// Si hubiera un ciclo de categorías (A dentro de B y B dentro de A), esas
	// categorías no cuelgan de ninguna raíz: las mostramos como principales.
	for _, c := range categorias {
		nodo := nodos[c.IDCategoria]
		if !visitados[nodo] {
			nodo.Hijas = nil
			calcularTotales(nodo, visitados)
			raices = append(raices, nodo)
		}
	}

	// Por la regla de negocio no debería haber productos sin categoría,
	// pero si los hubiera (datos viejos), no se pierden de vista.
	if len(sinCategoria.Productos) > 0 {
		calcularTotales(sinCategoria, visitados)
		raices = append(raices, sinCategoria)
	}

	deposito := Deposito{Raices: raices}
	for _, r := range raices {
		deposito.Unidades += r.Unidades
		deposito.Inversion += r.Inversion
	}
	return deposito
}

// calcularTotales suma las unidades e inversión del nodo y de todas sus subcategorías.
func calcularTotales(n *NodoCategoria, visitados map[*NodoCategoria]bool) {
	if visitados[n] {
		return
	}
	visitados[n] = true

	n.Unidades, n.Inversion = 0, 0
	for _, p := range n.Productos {
		n.Unidades += int64(p.Stock)
		n.Inversion += InversionDe(p)
	}
	for _, h := range n.Hijas {
		calcularTotales(h, visitados)
		n.Unidades += h.Unidades
		n.Inversion += h.Inversion
	}
}
