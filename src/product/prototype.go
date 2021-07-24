package product

import (
	"watcharis/ywd-test/model"

	"github.com/labstack/echo/v4"
)

type ProductService interface {
	AddProduct(ctx echo.Context) error
	GetProduct(ctx echo.Context) error
	FileProduct(ctx echo.Context) error
	TestConcerency(ctx echo.Context) error
}

type ProductRepository interface {
	InsertProduct(data model.ProductRequest) string
	FindProductByName(data model.ProductRequest) []model.Product
	FindAllProduct() []model.Product
}

var UrlPython = []string{
	"http://localhost:4567/nodepromise1",
	"http://127.0.0.1:4567/nodepromise2",
	"http://localhost:4567/nodepromise3",
	"http://localhost:4567/nodepromise4",
}

type Data struct {
	One   int `json:"one"`
	Two   int `json:"two"`
	Three int `json:"three"`
}
