package main

import (
	"github.com/braxes-backend/cmd"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", cmd.IndexHandler)
	http.HandleFunc("/orders", cmd.ListOrders)
	http.HandleFunc("/history", cmd.OrdersHistory)
	http.HandleFunc("/orders/{id}", cmd.OrderDetails)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
