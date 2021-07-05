package repository

import (
	"fmt"
	"time"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/receipt"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type Mysql struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) receipt.ReceiptRepository {
	return &Mysql{
		db: db,
	}
}

func (query *Mysql) FindSlipByReceiptCode(receiptCode string) ([]model.Receipt, error) {

	// sqlCommand := `SELECT * FROM receipt WHERE receipt_code=?`
	// query.db.Raw(sqlCommand, receiptCode).Scan(&receipts)
	var receipts []model.Receipt

	if err := query.db.Table("receipt").Where("receipt_code=?", receiptCode).Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

func (query *Mysql) InsertSlip(data model.ReceiptRequestBody) (string, error) {
	sqlCommand := `INSERT INTO receipt (receipt_id, user_id, receipt_code, receipt_point, total_price, status_receipt, create_date)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := query.db.Exec(sqlCommand, uuid.NewV4(), data.UserId, data.ReceiptCode, nil, data.TotalPrice, 0, time.Now()).Error; err != nil {
		return "", err
	}
	return "insert slip success", nil
}

func (query *Mysql) FindAllSlipByStatus() ([]model.Receipt, error) {
	// sqlCommand := `SELECT * FROM receipt WHERE status_receipt=?`
	var receipt []model.Receipt
	// query.db.Raw(sqlCommand, 0).Scan(&receipt)
	if err := query.db.Table("receipt").Where("status_receipt=?", 0).Find(&receipt).Error; err != nil {
		return nil, err
	}
	return receipt, nil
}

func (query *Mysql) FindSlipById(receiptId string) ([]model.Receipt, error) {
	sqlCommand := `SELECT * FROM receipt WHERE receipt_id=?`
	var receipt []model.Receipt
	if err := query.db.Raw(sqlCommand, receiptId).Scan(&receipt).Error; err != nil {
		return nil, err
	}
	return receipt, nil
}

func (query *Mysql) UpdateReceiptPoint(receiptId string, receiptPoint int16) (string, error) {
	sqlCommand := `UPDATE receipt SET receipt_point=?, status_receipt=1 WHERE receipt_id=?`
	if updatePoint := query.db.Exec(sqlCommand, receiptPoint, receiptId).Error; updatePoint != nil {
		fmt.Println("updatePoint :", updatePoint)
		return "", updatePoint
	}
	return "update receipt point success", nil
}

func (query *Mysql) SumReceiptPoint(userId string) (int16, error) {
	sqlCommand := `SELECT sum(receipt_point) FROM receipt WHERE user_id=?`
	var totalReceipt int16
	if err := query.db.Raw(sqlCommand, userId).Scan(&totalReceipt).Error; err != nil {
		// has error
		return 0, err
	}
	// no error
	return totalReceipt, nil
}
