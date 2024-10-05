package main

import (
	_ "embed"
	"os"

	"github.com/braxes-backend/app/handlers"
	"github.com/braxes-backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed database/orders/sql/schema.sql
var ordersDDL string

func main() {
	log.Info("Starting api...")
	log.Info(os.Getenv("SQLITE_DATA"))
	database.Connect()
	defer database.DB.Close()
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
