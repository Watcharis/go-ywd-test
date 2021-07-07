package controller

import (
	"watcharis/ywd-test/src/users"

	"github.com/labstack/echo/v4"
)

type userservice struct {
	service users.UserService
}

func NewUserHandle(service users.UserService) *userservice {
	return &userservice{
		service: service,
	}
}

func (handle userservice) RouteGroup(e *echo.Group) {

	router := e.Group("/api/v1/user")
	router.POST("/register", handle.service.Register)
	router.GET("/login", handle.service.Login)
	router.POST("/exchange", handle.service.Exchange)
	router.GET("/getexchange/:userid", handle.service.UserExchange)

}
