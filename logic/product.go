package logic

import "errors"

type Product struct {
	ID          int     `json:"id_producto"`
	Name        string  `json:"nombre"`
	Description *string `json:"descripcion,omitempty"`
	Price       float64 `json:"precio"`
	Stock       int     `json:"stock"`
	CategoryID  *int    `json:"id_categoria,omitempty"`
}

func ValidateProduct(p Product) error {
	if p.Name == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if p.Price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	if p.Stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}
	return nil
}

func ApplyDiscount(p Product, percentage float64) Product {
	p.Price = p.Price - (p.Price * percentage / 100)
	return p
}
