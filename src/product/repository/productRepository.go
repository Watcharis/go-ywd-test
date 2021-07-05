package repository

import (
	"time"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/product"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

// สร้าง struct methods เพื่อ implement interface
type Mysql struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.ProductRepository {
	return &Mysql{
		db: db,
	}
}

func (query *Mysql) InsertProduct(data model.ProductRequest) string {
	product_id := uuid.NewV4()
	sqlCommand := `INSERT INTO products (product_id, product_name, product_point, create_date) VALUES (?, ?, ?, ?)`
	query.db.Exec(sqlCommand, product_id, data.ProductName, data.ProductPoint, time.Now())
	return "insert product success"
}

func (query *Mysql) FindProductByName(data model.ProductRequest) []model.Product {
	sqlCommand := `SELECT * FROM products WHERE product_name=?`
	var products []model.Product
	query.db.Raw(sqlCommand, data.ProductName).Scan(&products)
	return products
}

func (query *Mysql) FindAllProduct() []model.Product {
	sqlCommand := `SELECT * FROM products`
	var products []model.Product
	query.db.Raw(sqlCommand).Scan(&products)
	return products
}
