package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func SetConnection() *sql.DB {
	connectionDB, err := sql.Open("mysql", "root:9a1b3m@tcp(127.0.0.1:3306)/api_class_rest")
	if err != nil {
		fmt.Println("Error to connect", err)
	}

	return connectionDB
}
