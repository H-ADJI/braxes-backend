package database

import "database/sql"

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
	DB, _ = sql.Open("sqlite3", ":memory:")
}
