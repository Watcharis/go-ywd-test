package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMysqlDB() (*gorm.DB, error) {

	var (
		DBUSER     string = os.Getenv("DBUSER")
		DBPASSWORD string = os.Getenv("DBPASSWORD")
		DBHOST     string = os.Getenv("DBHOST")
		DBPORT     string = os.Getenv("DBPORT")
		DBNAME     string = os.Getenv("DBNAME")
	)
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`, DBUSER, DBPASSWORD, DBHOST, DBPORT, DBNAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
