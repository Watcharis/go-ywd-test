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
