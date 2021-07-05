package model

import "time"

type ProductRequest struct {
	ProductName  string `json:"product_name"`
	ProductPoint string `json:"product_point"`
}

type Product struct {
	ProductId    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPoint string    `json:"product_point"`
	CreateDate   time.Time `json:"create_date"`
}
