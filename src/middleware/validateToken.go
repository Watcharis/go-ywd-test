package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"watcharis/ywd-test/database"
	"watcharis/ywd-test/model"
	claimsUsers "watcharis/ywd-test/src/users/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func AuthenCheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			header := ctx.Request().Header
			testHeader := ctx.Request().Header.Get("Accesstoken")

			fmt.Println("header :", testHeader)
			db, err := database.ConnectMysqlDB()
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			if accesstToken, existsKey := header["Accesstoken"]; existsKey {
				if token := strings.Split(accesstToken[0], " "); len(token) == 2 {

					// JWT Parse() จะทำการ validate token ที่เข้ามาให้
					validateToken, err := jwt.Parse(token[1], func(*jwt.Token) (interface{}, error) {
						return []byte("goywdtest"), nil
					})

					if err != nil {
						logrus.Errorln("err validate Token ->", err.Error())
						return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
					}

					//TODO validate token
					claims := validateToken.Claims

					//convert data to []byte
					tmp, _ := json.Marshal(claims)

					var tokenClaim claimsUsers.Claims

					//encode []byte to map, slice, array, ....
					if err := json.Unmarshal(tmp, &tokenClaim); err != nil {
						return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
					}
					// fmt.Println("tokenClaim :", tokenClaim)
					// fmt.Println("time :", time.Now().Unix())            ###time now

					var users []model.Users
					if err := db.Table("users").Where("email=?", tokenClaim.Email).Find(&users).Error; err != nil {
						return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
					}

					if len(users) != 0 {
						if tokenClaim.ExpiresAt > time.Now().Unix() {

							// TO DO set Data in middleware ctx.Set()
							ctx.Set("user_id", tokenClaim.UserId)
							ctx.Set("role_id", tokenClaim.RoleId)
							ctx.Set("email", tokenClaim.Email)
							return next(ctx)
						}
						return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "token expires", Status: "fail", Data: ""})
					}
					return ctx.JSON(http.StatusNotFound, model.JsonResponse{Message: "not found user", Status: "fail", Data: ""})
				} else {
					return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: "Unauthorized", Status: "fail", Data: ""})
				}
			}
			return ctx.JSON(http.StatusNotFound, model.JsonResponse{Message: "header undefine key accesstoken", Status: "fail", Data: ""})
		}
	}
}

func RejectRoleUnderAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// TO DO call Data from middleware ctx.Get()
			email := ctx.Get("email")
			roleId := ctx.Get("role_id")

			// call db
			db, err := database.ConnectMysqlDB()
			// fmt.Println("db :", db)

			if err != nil {
				return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			//TO DO check user in DB
			var users []model.Users
			if err := db.Table("users").Where("email=?", email).Find(&users).Error; err != nil {
				return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			if len(users) != 0 {
				if roleId == users[0].RoleId {
					// TO DO check Role user Admin
					var roles []model.Roles
					if err := db.Table("roles").Where("role_id=?", roleId).Find(&roles).Error; err != nil {
						return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
					}

					if roles[0].RoleName == "admin" {
						ctx.Set("user_id", users[0].UserId)
						ctx.Set("role_id", roles[0].RoleId)
						return next(ctx)
					}
					return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: "permission denine", Status: "fail", Data: ""})
				}
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "", Status: "fail", Data: ""})
			}
			return ctx.JSON(http.StatusNotFound, model.JsonResponse{Message: "not found user", Status: "fail", Data: ""})
		}
	}
}
