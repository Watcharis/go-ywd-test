package controller

import (
	"watcharis/ywd-test/src/product"

	"github.com/labstack/echo/v4"
)

type productService struct {
	service product.ProductService
}

func NewProductController(service product.ProductService) *productService {
	return &productService{
		service: service,
	}
}

func (h productService) RouteGroup(e *echo.Group) {
	router := e.Group("/api/v1/product")
	router.POST("", h.service.AddProduct)
	router.GET("", h.service.GetProduct)
	router.POST("/fileupload", h.service.FileProduct)
	router.GET("/testcon", h.service.TestConcerency)
}
