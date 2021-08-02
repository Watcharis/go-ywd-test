package main

import (
	"fmt"
	"net/http"

	database "watcharis/ywd-test/database"
	"watcharis/ywd-test/model"
	rdc "watcharis/ywd-test/redis"

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

	lineController "watcharis/ywd-test/src/line/controller"
	lineService "watcharis/ywd-test/src/line/service"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//TODO Echo register validate
	e.Validator = &CustomValidator{validator: validator.New()}
	e.IPExtractor = echo.ExtractIPDirect()

	//TODo CorsOrigin
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	//TODO GET ENV
	godotenv.Load(".env")

	//TODO connect Mysql DB
	db, err := database.ConnectMysqlDB()

	if err != nil {
		logrus.Errorln("error connect database ->", err.Error())
	}

	//TODO CREATE TABLE DB
	database.CreateDb(db)

	//TODO INIT PERMISSION USER IN DB
	database.InitPermissionUsersInDB(db)

	//TODO INIT ADMIN IN DB
	database.CreateAdminDB(db)

	//TODO GORM stage migrations
	// if err := migrate.MigrateDatabase(db, migrate.Table); err != nil {
	// 	logrus.Errorln("err migrations ->", err.Error())
	// }

	if db != nil {
		fmt.Println("db :", db)
	}

	//TODO Connect Redis
	rdb := rdc.NewConnectRedis()
	rct := rdb.ConnectRedis()
	// Api Example
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	routerPublich := e.Group("/publice")

	_userrepository := userRepository.NewUserRepository(db)
	_userservice := userService.NewUserService(_userrepository, rct)
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

	_lineService := lineService.NewLineService()
	_lineController := lineController.NewLineController(_lineService)
	_lineController.RouteGroup(routerPublich)

	e.Logger.Fatal(e.Start(":1323"))
}
