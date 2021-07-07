package main

import (
	"fmt"
	"net/http"

	database "watcharis/ywd-test/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

func main() {

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	godotenv.Load(".env")

	db, err := database.ConnectMysqlDB()

	if err != nil {
		panic(err.Error())
	}

	database.CreateDb(db)
	database.InitDatabase(db)

	if err != nil {
		return
	}

	if db != nil {
		fmt.Println("db :", *db)
	}

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
