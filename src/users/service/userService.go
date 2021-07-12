package service

import (
	"fmt"
	"net/http"
	"time"

	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/users"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
	Email  string `json:"email"`
	*jwt.StandardClaims
}

type Usersrepository struct {
	repository users.UserRepository
}

func NewUserService(repository users.UserRepository) *Usersrepository {
	return &Usersrepository{
		repository: repository,
	}
}

func (r *Usersrepository) Register(ctx echo.Context) error {

	bodyUser := users.UserRequestRegister{}

	if reqbody := ctx.Bind(&bodyUser); reqbody != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: reqbody.Error(), Status: "fall", Data: ""})
	}

	if err := ctx.Validate(&bodyUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fall", Data: ""})
	}

	fmt.Println("bodyUser :", bodyUser)
	checkQuery, err := r.repository.ValidateDataUserByEmail(bodyUser.Email)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if len(checkQuery) == 0 {

		getRoleId, err := r.repository.GetRoleIdByRoleName(bodyUser.RoleName)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		dataGenUser := model.UsersCreate{
			Email:       bodyUser.Email,
			Password:    bodyUser.Password,
			Phone:       bodyUser.Phone,
			DisplayName: bodyUser.DisplayName,
			RoleId:      getRoleId,
		}

		register, err := r.repository.RegisterUser(dataGenUser)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: register, Status: "success", Data: ""})
	}
	return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "User is exists", Status: "fail", Data: ""})
}

func (r *Usersrepository) Login(ctx echo.Context) error {

	bodyUser := users.UserRequestLogin{}

	if reqbody := ctx.Bind(&bodyUser); reqbody != nil {
		// fmt.Println("reqbody :", reqbody)
		return reqbody
	}

	if err := ctx.Validate(&bodyUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fall", Data: ""})
	}
	// fmt.Println("bodyUser :", bodyUser)

	data := users.UserRequestLogin{
		Email:    bodyUser.Email,
		Password: bodyUser.Password,
	}
	valiDateUser, err := r.repository.ValidateUserLogin(data)

	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if len(valiDateUser) != 0 {

		// fmt.Println("now :", time.Now())
		// fmt.Println("after 1 hour :", time.Now().Add(time.Hour*1))
		// gen token JWT

		claims := &Claims{
			valiDateUser[0].UserId,
			valiDateUser[0].RoleId,
			valiDateUser[0].Email,
			&jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// fmt.Println("token :", token)

		t, err := token.SignedString([]byte("goywdtest"))
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		// fmt.Println("t :", t)

		decodeToken, err := jwt.ParseWithClaims(t, claims, func(*jwt.Token) (interface{}, error) {
			return []byte("goywdtest"), nil
		})

		fmt.Printf("decodeToken : %v", *decodeToken)

		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "login success", Status: "success", Data: users.LoginResponse{Accesstoken: t}})
	} else {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "Not Found Users", Status: "fail", Data: ""})
	}
}
