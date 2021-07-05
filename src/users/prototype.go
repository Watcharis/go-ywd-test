package users

import (
	"watcharis/ywd-test/model"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
}

type UserRequestRegister struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	RoleName    string `json:"role_name"`
}

type UserRequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Accesstoken string `json:"accesstoken"`
}

type UserRepository interface {
	ValidateDataUserByEmail(email string) ([]model.Users, error)
	GetRoleIdByRoleName(roleName string) (string, error)
	InsertUser(data model.UsersCreate) (string, error)
	ValidateUserLogin(data UserRequestLogin) ([]model.Users, error)
}
