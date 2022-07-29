package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dsn := "root@tcp(localhost:3306)/course?"
	// dsn += "&charset=utf8"
	// dsn += "&interpolateParams=true"
	fmt.Println(dsn)
	if dsn == "" {
		log.Fatal("$DATABASE_URL is not set")
	}
	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
