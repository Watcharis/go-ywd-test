package product

import (
	"watcharis/ywd-test/model"

	"github.com/labstack/echo/v4"
)

type ProductService interface {
	AddProduct(ctx echo.Context) error
	GetProduct(ctx echo.Context) error
}

type ProductRepository interface {
	InsertProduct(data model.ProductRequest) string
	FindProductByName(data model.ProductRequest) []model.Product
	FindAllProduct() []model.Product
}
