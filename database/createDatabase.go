package database

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateDb(db *gorm.DB) {

	txRoles := db.Exec(
		`
			CREATE TABLE IF NOT EXISTS roles (
				role_id varchar(255) NOT NULL,
				role_name varchar(255) NOT NULL,
				role_status int NOT NULL,
				PRIMARY KEY(role_id)
			)
		`,
	)
	if txRoles.Error != nil {
		fmt.Printf("error : %s", txRoles.Error)
		panic(txRoles.Error)
	}

	txUser := db.Exec(
		`
			CREATE TABLE IF NOT EXISTS users (
				user_id varchar(255) NOT NULL ,
				email varchar(255) NOT NULL UNIQUE,
				password varchar(255) NOT NULL,
				phone varchar(255) NOT NULL,
				display_name varchar(255) NOT NULL,
				role_id varchar(255) NOT NULL,
				create_date DATETIME,
				PRIMARY KEY(user_id),
				CONSTRAINT FK_role_id FOREIGN KEY (role_id) REFERENCES roles(role_id)
			)
		`,
	)

	if txUser.Error != nil {
		fmt.Printf("error : %s", txUser.Error)
		panic(txUser.Error)
	}

	txProduct := db.Exec(
		`
		CREATE TABLE IF NOT EXISTS products (
			product_id varchar(255) NOT NULL,
			product_name varchar(255) NOT NULL,
			product_point int NOT NULL,
			create_date DATETIME,
			PRIMARY KEY(product_id)
		)
		`,
	)

	if txProduct.Error != nil {
		fmt.Printf("error : %s", txProduct.Error)
		panic(txProduct.Error)
	}

	txReceipt := db.Exec(
		`
		CREATE TABLE IF NOT EXISTS receipt (
			receipt_id varchar(255) NOT NULL ,
			user_id varchar(255) NOT NULL,
			receipt_code varchar(255) NOT NULL,
			receipt_point int NULL,
			total_price int NOT NULL,
			status_receipt int NOT NULL,
			create_date DATETIME,
			PRIMARY KEY(receipt_id),
			CONSTRAINT FK_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
		)
	`,
	)

	if txReceipt.Error != nil {
		fmt.Printf("error : %s", txReceipt.Error)
		panic(txReceipt.Error)
	}

	txUserProduct := db.Exec(
		`
		CREATE TABLE IF NOT EXISTS user_product (
			user_product_id varchar(255) NOT NULL ,
			user_id varchar(255) NOT NULL,
			product_id varchar(255) NOT NULL,
			create_date DATETIME,
			PRIMARY KEY(user_product_id),
			CONSTRAINT FK_user_pro_id FOREIGN KEY (user_id) REFERENCES users(user_id),
			CONSTRAINT FK_product_use_id FOREIGN KEY (product_id) REFERENCES products(product_id)
		)
		`,
	)

	if txUserProduct.Error != nil {
		fmt.Printf("error : %s", txUserProduct.Error)
		panic(txUserProduct.Error)
	}

}
