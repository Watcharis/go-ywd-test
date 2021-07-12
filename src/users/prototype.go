package users

import (
	"watcharis/ywd-test/model"

	"github.com/labstack/echo/v4"
)

type UserRequestRegister struct {
	Email       string `json:"email" validate:"required,emai"`
	Phone       string `json:"phone" validate:"required,phone"`
	Password    string `json:"password" validate:"required,password"`
	DisplayName string `json:"display_name" validate:"required,display_name"`
	RoleName    string `json:"role_name" validate:"required,role_name"`
}

type UserRequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Accesstoken string `json:"accesstoken"`
}

type UserAndProduct struct {
	UserId    string `json:"user_id"  validate:"required"`
	ProductId string `json:"product_id"`
}

type UserService interface {
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
	Exchange(ctx echo.Context) error
	UserExchange(ctx echo.Context) error
}

type UserRepository interface {
	ValidateDataUserByEmail(email string) ([]model.Users, error)
	GetRoleIdByRoleName(roleName string) (string, error)
	RegisterUser(data model.UsersCreate) (string, error)
	ValidateUserLogin(data UserRequestLogin) ([]model.Users, error)
	SumReceiptPoint(userId string) (int16, error)
	FindProductByid(productId string) ([]model.Product, error)
	FindUserProductByUserId(userId string) ([]model.UserJoinProducts, error)
	InsertUserProduct(data UserAndProduct) (string, error)
}
