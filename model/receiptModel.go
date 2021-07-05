package model

import "time"

type ReciptRequest struct {
	ReceiptId   string `json:"receipt_id"`
	UserId      string `json:"user_id"`
	ReceiptCode string `json:"receipt_code"`
	TotalPrice  string `json:"total_price"`
}

type Receipt struct {
	ReceiptId     string    `json:"receipt_id"`
	UserId        string    `json:"user_id"`
	ReceiptCode   string    `json:"receipt_code"`
	ReceiptPoint  int16     `json:"receipt_point"`
	TotalPrice    int16     `json:"total_price"`
	StatusReceipt int16     `json:"status_receipt"`
	CreateDate    time.Time `json:"create_date"`
}

type ReceiptRequestBody struct {
	UserId      string `json:"user_id" validate:"required,user_id"`
	ReceiptCode string `json:"receipt_code" validate:"required,receipt_code"`
	TotalPrice  string `json:"total_price" validate:"required,receipt_code"`
}

type ReceiptUpdatePointById struct {
	ReceiptId string `json:"receipt_id"`
}
