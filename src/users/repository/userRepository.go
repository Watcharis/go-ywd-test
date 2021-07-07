package repository

import (
	"time"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/users"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type Mysql struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.UserRepository {
	return &Mysql{
		db: db,
	}
}

func (query *Mysql) ValidateDataUserByEmail(email string) ([]model.Users, error) {
	var users []model.Users
	// sqlCommand := "SELECT * FROM users WHERE email=?"
	// query.db.Raw(sqlCommand, email).Scan(&users)
	if err := query.db.Table("users").Where("email=?", email).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (query *Mysql) GetRoleIdByRoleName(roleName string) (string, error) {
	sqlCommand := `SELECT roles.role_id FROM roles WHERE role_name=?`
	rows, err := query.db.Raw(sqlCommand, roleName).Rows()
	defer rows.Close()
	if err != nil {
		return "", err
	}
	var roleId string
	for rows.Next() {
		rows.Scan(&roleId)
	}
	return roleId, nil
}

func (query *Mysql) InsertUser(data model.UsersCreate) (string, error) {
	user_id := uuid.NewV4()
	sqlCommand := `INSERT INTO users (user_id, email, password, phone, display_name, role_id, create_date) VALUES (?, ?, ?, ?, ? ,?, ?)`
	if err := query.db.Exec(sqlCommand, user_id, data.Email, data.Password, data.Phone, data.DisplayName, data.RoleId, time.Now()).Error; err != nil {
		return "", err
	}
	return "insert users success", nil
}

func (query *Mysql) ValidateUserLogin(data users.UserRequestLogin) ([]model.Users, error) {
	sqlCommand := `SELECT * FROM users WHERE email=? AND password=?`
	var users []model.Users
	if err := query.db.Raw(sqlCommand, data.Email, data.Password).Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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

func (query *Mysql) FindProductByid(productId string) ([]model.Product, error) {
	var products []model.Product
	if err := query.db.Table("products").Where("product_id=?", productId).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (query *Mysql) FindUserProductByUserId(userId string) ([]model.UserJoinProducts, error) {
	sqlCommand := `
				SELECT * FROM user_product 
				LEFT JOIN products ON user_product.product_id = products.product_id
				WHERE user_product.user_id=?
				`
	var userProducts []model.UserJoinProducts
	if err := query.db.Raw(sqlCommand, userId).Find(&userProducts).Error; err != nil {
		return nil, err
	}
	// fmt.Println("userProducts :", userProducts)
	return userProducts, nil
}

func (query *Mysql) InsertUserProduct(data users.UserAndProduct) (string, error) {
	userProducrId := uuid.NewV4()
	sqlCommand := `INSERT INTO user_product (user_product_id, user_id, product_id, create_date) VALUES (?, ?, ?, ?)`
	if err := query.db.Exec(sqlCommand, userProducrId, data.UserId, data.ProductId, time.Now()).Error; err != nil {
		return "", err
	}
	return "insert user product success", nil
}
