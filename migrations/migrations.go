package migrations

import (
	"fmt"
	"watcharis/ywd-test/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Table = []interface{}{
	&model.Role{},
	&model.User{},
	&model.Products{},
	&model.Receipts{},
	&model.UserProducts{},
}

func MigrateDatabase(db *gorm.DB, tables []interface{}) error {
	tx := db.Begin()
	for _, t := range tables {
		fmt.Println("t :", t)
		if err := tx.AutoMigrate(t); err != nil {
			logrus.Errorln("err auto migration ->", err.Error())
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
