package cmd

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const dbName = "orders.db"

type Order struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	ImageUrl    string `json:"image_url"`
	Ordered_at  int    `json:"ordered_at"`
	Added_at    int    `json:"added_at"`
	Prepared_at int    `json:"prepared_at"`
	IsPrepared  bool   `json:"isPrepared"`
	Details     string `json:"details"`
	Date        string `json:"date"`
}
type OrderDB struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string
	ImageUrl    sql.NullString `json:"imageUrl"   db:"image_url"`
	Ordered_at  sql.NullInt64
	Added_at    sql.NullInt64
	Prepared_at int
	IsPrepared  bool `json:"isPrepared" db:"is_prepared"`
	Details     sql.NullString
}

func initDB() *sqlx.DB {

	db, err := sqlx.Open("sqlite", dbName)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	return db
}

var dbHandle *sqlx.DB

func init() {
	dbHandle = initDB()
	createTables()
}
func createTables() {
	q, err := dbHandle.Prepare(`CREATE TABLE IF NOT EXISTS orders 
(
    id INTEGER PRIMARY KEY,
    name TEXT, url TEXT, 
    ordered_at INTEGER DEFAULT 0, added_at INTEGER DEFAULT 0, prepared_at INTEGER DEFAULT 0,
    is_prepared INTEGER DEFAULT 0, image_url TEXT, details TEXT 
)`)
	defer q.Close()
	if err != nil {
		log.Fatal("Incorrect statement for create table")
	}
	_, err = q.Exec()

	if err != nil {
		log.Fatal("couldn't create orders table")
	}

}

func OrderById(id int) OrderDB {
	order := OrderDB{}
	err := dbHandle.Get(&order, "SELECT * FROM orders WHERE id= ?;", id)
	if err != nil {
		log.Fatal(err)
	}
	return order
}
func AllOrders() []OrderDB {
	orders := make([]OrderDB, 0)
	err := dbHandle.Select(&orders, "SELECT * FROM orders WHERE is_prepared=0;")
	if err != nil {
		log.Fatal(err)
	}
	return orders
}
func AllPreparedOrders() []OrderDB {
	orders := make([]OrderDB, 0)
	err := dbHandle.Select(&orders, "SELECT * FROM orders WHERE is_prepared=1;")
	if err != nil {
		log.Fatal(err)
	}
	return orders
}
