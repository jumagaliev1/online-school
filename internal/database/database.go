package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	// dsn := "root@tcp(localhost:3306)/course?"
	// dsn += "&charset=utf8"
	// dsn += "&interpolateParams=true"
	dsn := os.Getenv("CLEARDB_DATABASE_URL")
	dsn = strings.Replace(dsn, "mysql://", "", -1)
	dsn = strings.Replace(dsn, "@", "@tcp(", -1)
	dsn = strings.Replace(dsn, "net", "net)", -1)
	dsn = strings.Replace(dsn, "?reconnect=true", "", -1)

	fmt.Println(dsn)
	if dsn == "" {
		log.Fatal("$DATABASE_URL is not set")
	}

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
