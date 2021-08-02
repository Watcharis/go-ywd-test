package controller

import (
	"watcharis/ywd-test/src/line"

	"github.com/labstack/echo/v4"
)

type lineController struct {
	service line.LineService
}

func NewLineController(service line.LineService) lineController {
	return lineController{
		service: service,
	}
}

func (l lineController) RouteGroup(e *echo.Group) {
	router := e.Group("/api/v1/line")
	router.POST("/webhook", l.service.Webhook)
}
