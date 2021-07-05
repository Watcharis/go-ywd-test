package database

import (
	"errors"
	"fmt"
	"time"

	"watcharis/ywd-test/model"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

func getRoleIdByRoleName(db *gorm.DB) string {
	selectRoles := `SELECT roles.role_id FROM roles WHERE role_name=?`
	rows, err := db.Raw(selectRoles, "admin").Rows()
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var roles string
	for rows.Next() {
		rows.Scan(&roles)
	}
	return roles
}

func getRoleByRoleName(db *gorm.DB, permission string) string {
	selectRoles := `SELECT roles.role_name FROM roles WHERE role_name=?`
	rows, err := db.Raw(selectRoles, permission).Rows()
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var roles string
	for rows.Next() {
		rows.Scan(&roles)
	}
	return roles
}

func getAllRole(db *gorm.DB) []model.Roles {
	var roles []model.Roles
	selectRoles := `SELECT * FROM roles`
	db.Raw(selectRoles).Scan(&roles)
	return roles
}

func getUserAdminByEmail(db *gorm.DB) []model.Users {
	var users []model.Users
	selectUsers := `SELECT * FROM users WHERE email=?`
	db.Raw(selectUsers, "admin@admin.com").Scan(&users)
	return users
}

func InitDatabase(db *gorm.DB) {

	storePermission := []string{"admin", "user"}
	role_id := uuid.NewV4()
	status := 1

	for _, permission := range storePermission {

		checkRole := getRoleByRoleName(db, permission)

		if checkRole == "" {

			gennarateRoles := db.Exec(`INSERT IGNORE INTO roles (role_id, role_name, role_status) values(?, ?, ?)`, role_id, permission, status)

			if gennarateRoles.Error != nil {
				panic(gennarateRoles.Error)
			}

			status = status - 1
		} else {
			fmt.Println(errors.New("duplicate roles in DB"))
		}

	}

	getUsers := getUserAdminByEmail(db)

	if len(getUsers) > 0 {
		fmt.Println(errors.New("Users is exists"))
	} else {

		getRoleId := getRoleIdByRoleName(db)
		currentTime := time.Now()

		createAdmin := db.Exec(
			`INSERT IGNORE INTO users (user_id, email, password, phone, display_name, role_id, create_date) values(?, ?, ?, ?, ?, ?, ?)`,
			uuid.NewV4(), "admin@admin.com", "admin", "0999999999", "admin", getRoleId, currentTime,
		)

		if createAdmin.Error != nil {
			panic(createAdmin.Error)
		}
	}

}
