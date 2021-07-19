package controller

import (
	"watcharis/ywd-test/src/receipt"

	"watcharis/ywd-test/src/middleware"

	"github.com/labstack/echo/v4"
)

type receiptService struct {
	service receipt.ReceiptService
}

func NewReceiptContrller(service receipt.ReceiptService) *receiptService {
	return &receiptService{
		service: service,
	}
}

func (handle receiptService) RouteGroup(e *echo.Group) {

	router := e.Group("/api/v1/receipt")
	router.Use(middleware.AuthenCheckToken())
	router.POST("/sendslip", handle.service.SendSlip)
	router.GET("/getslip", handle.service.AdminGetSlip)
	router.PUT("/point", handle.service.AdminGiveReceiptPoint)
	router.GET("/totalpoint", handle.service.GetTotalPointUser)
}
