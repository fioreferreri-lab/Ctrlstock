package logic

import (
	"context"
	"fmt"
	"time"

	"ctrlstock/db"
)

// PeriodosVentas son los períodos (en días) que se pueden consultar en el ranking.
var PeriodosVentas = []int{7, 30, 90}

// DiasAlertaStock: si el stock alcanza para esta cantidad de días o menos, hay que reponer.
const DiasAlertaStock = 7

// ProductoVendido es una fila del ranking de más vendidos.
type ProductoVendido struct {
	IDProducto  int32
	Nombre      string
	Categoria   string
	Stock       int32
	Vendidas    int64
	DiasDeStock int  // para cuántos días alcanza el stock al ritmo de venta actual
	Alerta      bool // true si el stock alcanza para DiasAlertaStock días o menos
}

// RankingVentas es el resultado completo para la página de más vendidos.
type RankingVentas struct {
	Dias             int // período consultado
	Productos        []ProductoVendido
	TotalVendidas    int64
	MaxVendidas      int64 // las del primero, para dibujar las barras proporcionales
	ParaReponer      int   // cuántos productos están en alerta
	ProductoEstrella string
}

// periodoValido devuelve el período pedido si está permitido, o 30 días por defecto.
func periodoValido(dias int) int {
	for _, p := range PeriodosVentas {
		if p == dias {
			return dias
		}
	}
	return 30
}

// ObtenerRankingVentas arma el ranking de productos más vendidos en los últimos `dias` días.
// Solo cuenta egresos con motivo "Venta". Además calcula, al ritmo de venta del período,
// para cuántos días alcanza el stock de cada producto.
func ObtenerRankingVentas(ctx context.Context, queries *db.Queries, dias int) (RankingVentas, error) {
	dias = periodoValido(dias)
	desde := time.Now().AddDate(0, 0, -dias)

	filas, err := queries.RankingVentas(ctx, desde)
	if err != nil {
		return RankingVentas{}, fmt.Errorf("error al calcular el ranking de ventas: %w", err)
	}

	ranking := RankingVentas{Dias: dias}
	for _, f := range filas {
		porDia := float64(f.Vendidas) / float64(dias)
		diasDeStock := 0
		if porDia > 0 {
			diasDeStock = int(float64(f.Stock) / porDia)
		}

		item := ProductoVendido{
			IDProducto:  f.IDProducto,
			Nombre:      f.Nombre,
			Categoria:   f.Categoria,
			Stock:       f.Stock,
			Vendidas:    f.Vendidas,
			DiasDeStock: diasDeStock,
			Alerta:      diasDeStock <= DiasAlertaStock,
		}
		ranking.Productos = append(ranking.Productos, item)
		ranking.TotalVendidas += f.Vendidas
		if item.Alerta {
			ranking.ParaReponer++
		}
	}

	if len(ranking.Productos) > 0 {
		ranking.MaxVendidas = ranking.Productos[0].Vendidas
		ranking.ProductoEstrella = ranking.Productos[0].Nombre
	}
	return ranking, nil
}
