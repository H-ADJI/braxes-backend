package database

import (
	"database/sql"
	"os"
	"strconv"
)

// var DB DBInstance
//
// type DBInstance struct {
// 	Db *sql.DB
// }
//
// func Connect() {
// 	db, _ := sql.Open("sqlite3", ":memory:")
// 	DB = DBInstance{Db: db}
// }

var DB *sql.DB

func Connect() {
	isProd, _ := strconv.ParseBool(os.Getenv("IS_PROD"))
	if isProd {
		DB, _ = sql.Open("sqlite3", "/data/db.sqlite3")
		return
	}
	DB, _ = sql.Open("sqlite3", ":memory:")
}
