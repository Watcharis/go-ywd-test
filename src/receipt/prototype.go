package receipt

import (
	"watcharis/ywd-test/model"

	"github.com/labstack/echo/v4"
)

type ReceiptService interface {
	SendSlip(ctx echo.Context) error
	AdminGetSlip(ctx echo.Context) error
	AdminGiveReceiptPoint(ctx echo.Context) error
	GetTotalPointUser(ctx echo.Context) error
}

type ReceiptRepository interface {
	FindSlipByReceiptCode(receiptCode string) ([]model.Receipt, error)
	InsertSlip(data model.ReceiptRequestBody) (string, error)
	FindAllSlipByStatus() ([]model.Receipt, error)
	FindSlipById(receiptId string) ([]model.Receipt, error)
	UpdateReceiptPoint(receiptId string, receiptPoint int16) (string, error)
	SumReceiptPoint(userId string) (int16, error)
}
