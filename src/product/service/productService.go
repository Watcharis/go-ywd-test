package service

import (
	"net/http"
	"strconv"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/product"

	"github.com/labstack/echo/v4"
)

type productRepository struct {
	repository product.ProductRepository
}

func NewProductService(repository product.ProductRepository) *productRepository {
	return &productRepository{
		repository: repository,
	}
}

func (r productRepository) AddProduct(ctx echo.Context) error {
	bodyProduct := new(model.ProductRequest)

	if reqbody := ctx.Bind(bodyProduct); reqbody != nil {
		// fmt.Println("reqbody :", reqbody)
		return reqbody
	}

	if bodyProduct.ProductName == "" {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "invalide format ProductName", Status: "fail", Data: ""})
	}

	proDuctPoint, err := strconv.Atoi(bodyProduct.ProductPoint)
	if proDuctPoint <= 0 {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "ProductPonit must more than zero", Status: "fail", Data: ""})
	}

	if err != nil {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if bodyProduct.ProductName != "" && proDuctPoint > 0 {

		dataProduct := model.ProductRequest{
			ProductName:  bodyProduct.ProductName,
			ProductPoint: bodyProduct.ProductPoint,
		}

		validateProduct := r.repository.FindProductByName(dataProduct)

		if len(validateProduct) != 0 {
			return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "product is exists", Status: "fail", Data: ""})
		} else {
			insertItemProduct := r.repository.InsertProduct(dataProduct)

			if insertItemProduct == "insert product success" {
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: insertItemProduct, Status: "success", Data: ""})
			} else {
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "add product fail", Status: "fail", Data: ""})
			}
		}
	} else {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "product_point must more than zero and invalide format product_name", Status: "fail", Data: ""})
	}
}

func (r productRepository) GetProduct(ctx echo.Context) error {
	getProduct := r.repository.FindAllProduct()
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "get product success", Status: "success", Data: getProduct})
}
