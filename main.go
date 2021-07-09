package main

import (
	"fmt"
	"net/http"

	database "watcharis/ywd-test/database"
	"watcharis/ywd-test/model"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	userController "watcharis/ywd-test/src/users/controller"
	userRepository "watcharis/ywd-test/src/users/repository"
	userService "watcharis/ywd-test/src/users/service"

	productController "watcharis/ywd-test/src/product/controller"
	productRepository "watcharis/ywd-test/src/product/repository"
	productService "watcharis/ywd-test/src/product/service"

	receiptController "watcharis/ywd-test/src/receipt/controller"
	receiptRepository "watcharis/ywd-test/src/receipt/repository"
	receiptService "watcharis/ywd-test/src/receipt/service"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	return cv.validator.Struct(i)
}

func main() {

	e := echo.New()

	// //TODO Echo validate
	e.Validator = &CustomValidator{validator: validator.New()}

	//TODo CorsOrigin
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	//TODO GET ENV
	godotenv.Load(".env")
	db, err := database.ConnectMysqlDB()

	if err != nil {
		logrus.Errorln("error connect database ->", err.Error())
	}

	database.CreateDb(db)
	database.InitDatabase(db)

	if db != nil {
		fmt.Println("db :", *db)
	}

	// Api sample
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	routerPublich := e.Group("/publice")

	_userrepository := userRepository.NewUserRepository(db)
	_userservice := userService.NewUserService(_userrepository)
	_userController := userController.NewUserHandle(_userservice)
	_userController.RouteGroup(routerPublich)

	_productrepository := productRepository.NewProductRepository(db)
	_productservice := productService.NewProductService(_productrepository)
	_productcontroller := productController.NewProductController(_productservice)
	_productcontroller.RouteGroup(routerPublich)

	_receiptRepository := receiptRepository.NewReceiptRepository(db)
	_receiptservice := receiptService.NewReceiptService(_receiptRepository)
	_receiptcontroller := receiptController.NewReceiptContrller(_receiptservice)
	_receiptcontroller.RouteGroup(routerPublich)

	e.Logger.Fatal(e.Start(":1323"))
}
