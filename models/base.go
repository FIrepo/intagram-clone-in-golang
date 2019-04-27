package models

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
)

func PgCon() *sqlx.DB {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbUri)

	conn, err := sqlx.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}
	return conn
	// db.Debug().AutoMigrate(&Account{}, &Contact{})
}
