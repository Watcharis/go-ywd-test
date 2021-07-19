package service

import (
	"net/http"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/receipt"

	"github.com/labstack/echo/v4"
)

type receiptRepository struct {
	repository receipt.ReceiptRepository
}

func NewReceiptService(repository receipt.ReceiptRepository) *receiptRepository {
	return &receiptRepository{
		repository: repository,
	}
}

func (r *receiptRepository) SendSlip(ctx echo.Context) error {
	var bodyReceipt model.ReceiptRequestBody
	// bind data
	if err := ctx.Bind(bodyReceipt); err != nil {
		return err
	}

	//validate data
	if err := ctx.Validate(&bodyReceipt); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fall", Data: ""})
	}

	// fmt.Println("bodyReceipt :", bodyReceipt)

	validateReceiptCode, err := r.repository.FindSlipByReceiptCode(bodyReceipt.ReceiptCode)
	// fmt.Println("validateReceiptCode :", validateReceiptCode)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if len(validateReceiptCode) == 0 {

		data := &model.ReceiptRequestBody{
			UserId:      bodyReceipt.UserId,
			ReceiptCode: bodyReceipt.ReceiptCode,
			TotalPrice:  bodyReceipt.TotalPrice,
		}

		sendSlip, err := r.repository.InsertSlip(*data)

		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: sendSlip, Status: "success", Data: ""})
	}
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "receipt code is exists", Status: "fail", Data: ""})
}

func (r *receiptRepository) AdminGetSlip(ctx echo.Context) error {
	getSlip, err := r.repository.FindAllSlipByStatus()
	// fmt.Println("getSlip :", getSlip)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "get slip success", Status: "success", Data: getSlip})
}

func (r *receiptRepository) AdminGiveReceiptPoint(ctx echo.Context) error {

	// การ req Param ใน echo
	// receiptId := ctx.QueryParam("receipt_id")

	if valueReceipyId := ctx.QueryParam("receipt_id"); valueReceipyId != "" {

		checkReceiptId, err := r.repository.FindSlipById(valueReceipyId)

		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		if len(checkReceiptId) != 0 {

			if checkReceiptId[0].TotalPrice < 100 {
				giveReceiptPoint, err := r.repository.UpdateReceiptPoint(valueReceipyId, 30)
				if err != nil {
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: giveReceiptPoint, Status: "success", Data: ""})
			}

			var calTotalPrice int16
			calTotalPrice = checkReceiptId[0].TotalPrice % 100

			var point int16
			if calTotalPrice > 0 {

				point = ((checkReceiptId[0].TotalPrice/100)*100 + 30)

				giveReceiptPoint, err := r.repository.UpdateReceiptPoint(valueReceipyId, point)
				if err != nil {
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: giveReceiptPoint, Status: "success", Data: ""})
			}

			point = (checkReceiptId[0].TotalPrice / 100) * 100
			giveReceiptPoint, err := r.repository.UpdateReceiptPoint(valueReceipyId, point)
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}
			return ctx.JSON(http.StatusOK, model.JsonResponse{Message: giveReceiptPoint, Status: "success", Data: ""})
		}
	}
	return ctx.JSON(http.StatusNotFound, model.JsonResponse{Message: "receipt_id not found", Status: "fail", Data: ""})
}

func (r *receiptRepository) GetTotalPointUser(ctx echo.Context) error {
	if userId := ctx.QueryParam("user_id"); userId != "" {
		callSum, err := r.repository.SumReceiptPoint(userId)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "sum point success", Status: "success", Data: callSum})
	}
	return ctx.JSON(http.StatusNotFound, model.JsonResponse{Message: "not found user_id"})
}
