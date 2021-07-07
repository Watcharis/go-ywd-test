package database

import (
	"errors"
	"fmt"
	"time"

	"watcharis/ywd-test/model"

	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

func getRoleIdByRoleName(db *gorm.DB) (string, error) {
	selectRoles := `SELECT roles.role_id FROM roles WHERE role_name=?`
	rows, err := db.Raw(selectRoles, "admin").Rows()
	defer rows.Close()
	if err != nil {
		logrus.Errorln(err.Error())
	}
	var roles string
	for rows.Next() {
		rows.Scan(&roles)
	}
	return roles, nil
}

func getRoleByRoleName(db *gorm.DB, permission string) (string, error) {
	selectRoles := `SELECT roles.role_name FROM roles WHERE role_name=?`
	rows, err := db.Raw(selectRoles, permission).Rows()
	defer rows.Close()
	if err != nil {
		logrus.Errorln(err.Error())
	}
	var roles string
	for rows.Next() {
		rows.Scan(&roles)
	}
	return roles, nil
}

func getAllRole(db *gorm.DB) ([]model.Roles, error) {
	var roles []model.Roles
	selectRoles := `SELECT * FROM roles`
	if err := db.Raw(selectRoles).Scan(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func getUserAdminByEmail(db *gorm.DB) ([]model.Users, error) {
	var users []model.Users
	selectUsers := `SELECT * FROM users WHERE email=?`
	if err := db.Raw(selectUsers, "admin@admin.com").Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func InitDatabase(db *gorm.DB) {

	storePermission := []string{"admin", "user"}
	role_id := uuid.NewV4()
	status := 1
	tx := db.Begin()

	for _, permission := range storePermission {

		checkRole, err := getRoleByRoleName(db, permission)

		if err != nil {
			logrus.Errorln("getRoleByRoleName ->", err.Error())
		}

		if checkRole == "" {

			gennarateRoles := tx.Exec(`INSERT IGNORE INTO roles (role_id, role_name, role_status) values(?, ?, ?)`, role_id, permission, status)

			if gennarateRoles.Error != nil {
				defer logrus.Errorln("generate error ->", gennarateRoles.Error)
				tx.Rollback()
				// panic(gennarateRoles.Error)
			}

			status = status - 1
		} else {
			fmt.Println(errors.New("duplicate roles in DB"))
		}
	}

	getUsers, err := getUserAdminByEmail(db)

	if err != nil {
		logrus.Errorln("getUserAdminByEmail ->", err.Error())
	}

	if len(getUsers) > 0 {
		fmt.Println(errors.New("Users is exists"))
	} else {

		getRoleId, err := getRoleIdByRoleName(db)

		if err != nil {
			logrus.Errorln("getRoleByRoleName ->", err.Error())
		}
		currentTime := time.Now()

		createAdmin := db.Exec(
			`INSERT IGNORE INTO users (user_id, email, password, phone, display_name, role_id, create_date) values(?, ?, ?, ?, ?, ?, ?)`,
			uuid.NewV4(), "admin@admin.com", "admin", "0999999999", "admin", getRoleId, currentTime,
		)

		if createAdmin.Error != nil {
			defer logrus.Errorln("createAdmin error ->", createAdmin.Error)
			tx.Rollback()
			// panic(createAdmin.Error)
		}
	}

	tx.Commit()
}
