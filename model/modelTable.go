package model

import (
	"time"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type Model struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	*Model
	UserID      uuid.UUID `json:"user_id" gorm:"type:varchar(255);primaryKey;not null"`
	Email       string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password    string    `json:"password" gorm:"type:varchar(255);not null"`
	Phone       string    `json:"phone" gorm:"type:varchar(255);not null"`
	DisplayName string    `json:"display_name" gorm:"type:varchar(255);not null"`
	RoleID      uuid.UUID `json:"roleid" gorm:"type:varchar(255);not null"`
	Roles       Roles     `gorm:"foreignkey:RoleID;references:role_id"`
}

type Role struct {
	*Model
	RoleID     uuid.UUID `json:"role_id" gorm:"type:varchar(255);primaryKey;not null"`
	RoleName   string    `json:"role_name" gorm:"type:varchar(255);not null"`
	RoleStatus string    `json:"role_status" gorm:"type:varchar(10);not null"`
}

type Products struct {
	*Model
	ProductId    uuid.UUID `json:"product_id" gorm:"type:varchar(255);primaryKey"`
	ProductName  string    `json:"product_name" gorm:"type:varchar(255);not null"`
	ProductPoint string    `json:"product_point" gorm:"type:varchar(255);not null"`
}

type Receipts struct {
	*Model
	ReceiptId     uuid.UUID `json:"receipt_id" gorm:"type:varchar(255);primaryKey"`
	UserId        uuid.UUID `json:"user_id" gorm:"type:varchar(255) not null"`
	Users         Users     `gorm:"foreignkey:UserId;references:user_id"`
	ReceiptCode   string    `json:"receipt_code" gorm:"type:varchar(255) not null"`
	ReceiptPoint  int16     `json:"receipt_point" gorm:"type:int not null"`
	TotalPrice    int16     `json:"total_price" gorm:"type:int not null"`
	StatusReceipt string    `json:"status_receipt" gorm:"type:varchar(10);not null"`
}

type UserProducts struct {
	*Model
	UserProductId uuid.UUID `json:"user_product_id" gorm:"type:varchar(255);primaryKey"`
	UserId        uuid.UUID `json:"user_id" gorm:"type:varchar(255) not null"`
	Users         Users     `gorm:"foreignkey:UserId;references:user_id"`
	ProductId     uuid.UUID `json:"product_id" gorm:"type:varchar(255) not null"`
	Products      Products  `gorm:"foreignkey:ProductId;references:product_id"`
}
