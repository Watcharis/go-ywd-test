package model

import "time"

type UserRequestRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Display  string `json:"display"`
}

type JsonResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type Roles struct {
	RoleId     string `json:"role_id"`
	RoleName   string `json:"role_name"`
	RoleStatus int    `json:"role_status"`
}

type Users struct {
	UserId      string    `json:"user_id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Phone       string    `json:"phone"`
	DisplayName string    `json:"display_name"`
	RoleId      string    `json:"role_id"`
	CreateDate  time.Time `json:"create_date"`
}

type UsersCreate struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	DisplayName string `json:"display_name"`
	RoleId      string `json:"role_id"`
}
