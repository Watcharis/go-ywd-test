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
	bodyExchange := users.UserAndProduct{}
	// bodyExchange := new(users.UserAndProduct)

	if reqbody := ctx.Bind(&bodyExchange); reqbody != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: reqbody.Error(), Status: "fall", Data: ""})
	}
	if err := ctx.Validate(&bodyExchange); err != nil {
		return err
	}

	// fmt.Println("bodyExchange :", bodyExchange)
	dataExchange := users.UserAndProduct{
		UserId:    bodyExchange.UserId,
		ProductId: bodyExchange.ProductId,
	}

	//TODO query DB
	totalPoint, err := r.repository.SumReceiptPoint(dataExchange.UserId)
	productPoint, err := r.repository.FindProductByid(dataExchange.ProductId)
	checkExchange, err := r.repository.FindUserProductByUserId(dataExchange.UserId)

	//TODO convert string to int16
	resultProductPoint, err := strconv.Atoi(productPoint[0].ProductPoint)
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
			//คำนวณ เเต้มทั้งหมด
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: "point not enought exchange fail", Status: "fail", Data: totalPoint})
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

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "product not exists",
		"status":  "fail",
		"data":    "",
	})
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
