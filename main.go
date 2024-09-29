package main

import (
	_ "embed"

	"github.com/braxes-backend/app/handlers"
	"github.com/braxes-backend/database"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed database/orders/sql/schema.sql
var ordersDDL string

func main() {
	database.Connect()
	database.DB.Exec(ordersDDL)
	ordersRoutes := handlers.InitOrderHanlders()
	app := fiber.New(fiber.Config{
		ServerHeader: "braxes-backend",
		AppName:      "braxesApiv0.0.1",
	})
	// mount sub-routers
	app.Mount("/orders", ordersRoutes)
	app.Listen(":3000")

}
