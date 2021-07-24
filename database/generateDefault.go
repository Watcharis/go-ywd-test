package database

import (
	"errors"
	"time"

	"watcharis/ywd-test/model"

	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

func getRoleIdByRoleName(db *gorm.DB) ([]model.Roles, error) {
	tx := db.Begin()
	var roles []model.Roles
	if err := tx.Table("roles").Where("role_name=?", "admin").Find(&roles).Error; err != nil {
		return nil, err
	}
	tx.Commit()
	return roles, nil
}

func getRoleByRoleName(db *gorm.DB, permission string) ([]model.Roles, error) {
	tx := db.Begin()
	var roles []model.Roles
	if err := tx.Table("roles").Where("role_name=?", permission).Find(&roles).Error; err != nil {
		return nil, err
	}
	tx.Commit()
	return roles, nil
}

func getUserAdminByEmail(db *gorm.DB) ([]model.Users, error) {
	tx := db.Begin()
	var users []model.Users
	selectUsers := `SELECT * FROM users WHERE email=?`
	if err := tx.Raw(selectUsers, "admin@admin.com").Scan(&users).Error; err != nil {
		return nil, err
	}
	tx.Commit()
	return users, nil
}

func InitPermissionUsersInDB(db *gorm.DB) {
	tx := db.Begin()
	storePermission := []string{"admin", "user"}
	status := 1

	for _, permission := range storePermission {
		checkRole, err := getRoleByRoleName(db, permission)
		if err != nil {
			logrus.Errorln("getRoleByRoleName ->", err.Error())
		}

		if len(checkRole) == 0 {
			gennarateRoles := tx.Exec(`INSERT IGNORE INTO roles (role_id, role_name, role_status, create_date) values(?, ?, ?, ?)`, uuid.NewV4(), permission, status, time.Now())
			if gennarateRoles.Error != nil {
				logrus.Errorln("generate error ->", gennarateRoles.Error)
				tx.Rollback()
				// panic(gennarateRoles.Error)
			}
			status = status - 1
		} else {
			logrus.Warn("warning db table roles -> duplicate roles in DB")
		}
	}
	tx.Commit()
}

func CreateAdminDB(db *gorm.DB) {
	tx := db.Begin()

	//TODO Insert Admin
	getUsers, err := getUserAdminByEmail(db)
	if err != nil {
		logrus.Errorln("getUserAdminByEmail ->", err.Error())
	}

	if len(getUsers) > 0 {
		logrus.Errorln(errors.New("Users is exists"))
	} else {
		getRoleId, err := getRoleIdByRoleName(db)
		if err != nil {
			logrus.Errorln("getRoleByRoleName ->", err.Error())
		}
		createAdmin := tx.Exec(
			`INSERT IGNORE INTO users (user_id, email, password, phone, display_name, role_id, create_date) values(?, ?, ?, ?, ?, ?, ?)`,
			uuid.NewV4(), "admin@admin.com", "admin", "0999999999", "admin", getRoleId[0].RoleId, time.Now(),
		)
		if createAdmin.Error != nil {
			logrus.Errorln("createAdmin error ->", createAdmin.Error)
			tx.Rollback()
			// panic(createAdmin.Error)
		}
	}
	tx.Commit()
}
