package service

import (
	"net/http"
	"strconv"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/users"

	"github.com/labstack/echo/v4"
)

func calculateProduct(dataProduct []model.UserJoinProducts) int16 {
	var sum int16
	for i := 0; i < len(dataProduct); i++ {
		point, _ := strconv.Atoi(dataProduct[i].ProductPoint)
		sum += int16(point)
	}
	return sum
}

func (r *Usersrepository) Exchange(ctx echo.Context) error {

	bodyExchange := new(users.UserAndProduct)
	// fmt.Println("bodyExchange :", bodyExchange)

	if reqbody := ctx.Bind(&bodyExchange); reqbody != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: reqbody.Error(), Status: "fall", Data: ""})
	}
	// fmt.Println("bodyExchange :", bodyExchange)

	dataExchange := users.UserAndProduct{
		UserId:    bodyExchange.UserId,
		ProductId: bodyExchange.ProductId,
	}

	totalPoint, err := r.repository.SumReceiptPoint(dataExchange.UserId)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	productPoint, err := r.repository.FindProductByid(dataExchange.ProductId)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	checkExchange, err := r.repository.FindUserProductByUserId(dataExchange.UserId)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	// fmt.Println("totalPoint :", totalPoint)
	// fmt.Println("productPoint :", productPoint[0].ProductId)
	// fmt.Println("checkExchange :", checkExchange)

	resultProductPoint, err := strconv.Atoi(productPoint[0].ProductPoint)
	// fmt.Println("resultProductPoint :", resultProductPoint)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if resultProductPoint > 0 {
		if len(checkExchange) == 0 {

			if totalPoint > int16(resultProductPoint) {
				saveUserProduct, err := r.repository.InsertUserProduct(dataExchange)

				if err != nil {
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}

				if saveUserProduct == "insert user product success" {
					return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "exchange success", Status: "success", Data: ""})
				}

			}
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: "point not enought exchange fail", Status: "fail", Data: ""})
		} else {
			sumAllPointProductExchange := calculateProduct(checkExchange)
			// fmt.Println("sumAllPointProductExchange :", sumAllPointProductExchange)

			if (totalPoint - sumAllPointProductExchange) > int16(resultProductPoint) {

				saveUserProduct, err := r.repository.InsertUserProduct(dataExchange)
				if err != nil {
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}

				if saveUserProduct == "insert user product success" {
					return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "exchange success", Status: "success", Data: ""})
				}
			}
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: "point not enought exchange fail", Status: "fail", Data: ""})
		}
	}

	response := map[string]interface{}{
		"message": "product not exists",
		"status":  "fail",
		"data":    "",
	}

	return ctx.JSON(http.StatusOK, response)
}

func (r *Usersrepository) UserExchange(ctx echo.Context) error {
	userId := ctx.Param("userid")

	getProductAllOfUserId, err := r.repository.FindUserProductByUserId(userId)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	var resultAllProductOfUser []map[string]string

	for _, v := range getProductAllOfUserId {

		// productPoint, _ := strconv.Atoi(v.ProductPoint)
		// fmt.Println("productPoint :", productPoint)

		mapProductAllOfUser := map[string]string{
			"user_product_id": v.UserProductId,
			"product_id":      v.ProductId,
			"product_name":    v.ProductName,
			"product_point":   v.ProductPoint,
		}
		resultAllProductOfUser = append(resultAllProductOfUser, mapProductAllOfUser)
	}

	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "user exchange success", Status: "success", Data: resultAllProductOfUser})
}
