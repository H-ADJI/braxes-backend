package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates map[string]*template.Template
var templateNames = [...]string{"index", "table_row"}

func init() {
	templates = make(map[string]*template.Template)
	for _, templateName := range templateNames {
		templates[templateName] = template.Must(template.ParseFiles(
			fmt.Sprintf("templates/%s.html", templateName),
		))
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := templates["index"]
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
func ListOrders(w http.ResponseWriter, r *http.Request) {
	tmpl := templates["table_row"]
	orders := AllOrders()
	err := tmpl.Execute(w, orders)
	if err != nil {
		log.Fatal(err)
	}
}
func OrdersHistory(w http.ResponseWriter, r *http.Request) {
	tmpl := templates["table_row"]
	history := AllPreparedOrders()
	err := tmpl.Execute(w, history)
	if err != nil {
		log.Fatal(err)
	}
}
func OrderDetails(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	order := OrderById(id)
	jsonResponse, err := json.Marshal(order)
	if err != nil {
		// If there is an error in marshaling, return a server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set the HTTP status code to 200
	w.Write(jsonResponse)
}
func PrepareOrder(w http.ResponseWriter, r *http.Request) {
}
