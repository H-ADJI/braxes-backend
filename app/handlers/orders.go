package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/braxes-backend/app/middleware"
	"github.com/braxes-backend/database"
	"github.com/braxes-backend/database/orders"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// var queries *orders.Queries
const wixAPI = "https://www.wixapis.com/ecom/v1/orders/"

type reqOrder struct {
	PlatformID string
}
type RespOrder struct {
	Id           int64   `json:"id"`
	PlatformID   string  `json:"platformID"`
	TotalPrice   float64 `json:"totalPrice"`
	CustomerName string  `json:"CustomerName"`
	OrderNumber  int64   `json:"orderNumber"`
	CreationDate int64   `json:"creationDate"`
	IsProcessed  int64   `json:"isProcessed"`
}
type _ orders.Order

type shipping struct {
	deliveryPrice   string
	deliveryAddress string
}
type customer struct {
	name  string
	phone string
}
type item struct {
	title    string
	imageUrl string
	quantity float64
	price    string
}

func InitOrderHanlders() *fiber.App {
	orderQueries := orders.New(database.DB)
	dbCTX := context.Background()
	handler := fiber.New()
	handler.Use(middleware.InjectDB(orderQueries, dbCTX))
	handler.Get("", ListOrders)
	handler.Post("", AddOrder)
	handler.Get("/history", OrdersHistory)
	handler.Get("/unprocessed", ListUnprocessed)
	handler.Get("/processed", ListProcessed)
	handler.Get("/unprocess/:id", UnProcessOrder)
	handler.Get("/process/:id", ProcessOrder)
	handler.Post("/placed", OrderPlacedEvent)
	handler.Get("/:id", OrderDetails)
	return handler
}
func ListOrders(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	orders, err := q.GetAllOrders(dbCTX)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]RespOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			RespOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
				CustomerName: o.CustomerName,
				TotalPrice:   o.TotalPrice,
			},
		)
	}
	json, err := json.Marshal(ordersResp)
	if err != nil {
		log.Errorf("json marshal failed, %e \n", err)
		return fiber.ErrInternalServerError
	}
	return c.Send(json)
}
func ListUnprocessed(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	queries := c.Queries()
	isJson := queries["json"]
	orders, err := q.GetUnProcessedOrders(dbCTX)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]RespOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			RespOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				OrderNumber:  o.OrderNumber,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
				TotalPrice:   o.TotalPrice,
				CustomerName: o.CustomerName,
			},
		)
	}
	if isJson == "1" {
		json, err := json.Marshal(ordersResp)
		if err != nil {
			log.Errorf("json marshal failed, %e \n", err)
			return fiber.ErrInternalServerError
		}
		return c.Send(json)
	}
	trs := ordersTableRows(ordersResp)
	trCTX := context.Background()
	var bf bytes.Buffer
	trs.Render(trCTX, &bf)
	return c.Type("html").Send(bf.Bytes())
}
func ListProcessed(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	queries := c.Queries()
	isJson := queries["json"]
	orders, err := q.GetProcessedOrders(dbCTX)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]RespOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			RespOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
			},
		)
	}
	if isJson == "1" {
		json, err := json.Marshal(ordersResp)
		if err != nil {
			log.Errorf("json marshal failed, %e \n", err)
			return fiber.ErrInternalServerError
		}
		return c.Send(json)
	}
	trs := ordersTableRows(ordersResp)
	trCTX := context.Background()
	var bf bytes.Buffer
	trs.Render(trCTX, &bf)
	return c.Type("html").Send(bf.Bytes())
}
func AddOrder(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	reqBody := new(reqOrder)
	err := c.BodyParser(reqBody)
	if err != nil {
		log.Error(err)
		return fiber.ErrBadRequest
	}
	order, err := q.AddOrder(dbCTX, orders.AddOrderParams{
		PlatformID: reqBody.PlatformID, CreationDate: time.Now().Unix(),
	})
	if err != nil {
		log.Error(err)
		return fiber.NewError(fiber.StatusBadRequest, "order already exists on DB")
	}
	json, _ := json.Marshal(
		RespOrder{
			Id:           order.ID,
			PlatformID:   order.PlatformID,
			CreationDate: order.CreationDate,
			IsProcessed:  0,
		},
	)
	return c.Send(json)
}
func OrdersHistory(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	orders, err := q.GetAllOrdersDescDate(dbCTX)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]RespOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			RespOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
			},
		)
	}
	json, err := json.Marshal(ordersResp)
	if err != nil {
		log.Errorf("json marshal failed, %e\n", err)
		return fiber.ErrInternalServerError
	}
	return c.Send(json)
}
func ProcessOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	q, dbCTX := middleware.GetQueryCTX(c)
	processed, err := q.ProcessOrder(
		dbCTX,
		orders.ProcessOrderParams{
			ProcessedDate: sql.NullInt64{Valid: true, Int64: time.Now().Unix()}, PlatformID: id,
		},
	)
	if err != nil {
		log.Error(err)
		return fiber.ErrNotFound
	}
	if processed == 0 {
		log.Error("Invalid order state")
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}
func UnProcessOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	q, dbCTX := middleware.GetQueryCTX(c)
	processed, err := q.UnProcessOrder(
		dbCTX,
		id,
	)
	if err != nil {
		log.Error(err)
		return fiber.ErrNotFound
	}
	if processed == 1 {
		log.Error("Invalid order state")
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}
func OrderDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	// wixKey := os.Getenv("WIX_API_KEY")
	wixKey := os.Getenv("WIX_API_KEY")
	wixSiteId := os.Getenv("SITE_ID")
	wixAccountId := os.Getenv("ACCOUNT_ID")
	// TODO: add certificates in scratch image
	client := &http.Client{}
	req, err := http.NewRequest("GET", wixAPI+id, nil)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	req.Header.Add("wix-account-id", wixAccountId)
	req.Header.Add("wix-site-id", wixSiteId)
	req.Header.Add("Authorization", wixKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}
	var ordersDetails map[string]interface{}
	err = json.Unmarshal(content, &ordersDetails)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}
	queries := c.Queries()
	isJson := queries["json"]
	if isJson == "1" {
		return c.Status(fiber.StatusOK).JSON(ordersDetails)
	}
	order, ok := ordersDetails["order"].(map[string]interface{})
	if ok == false {
		log.Error("could not access order value")
		return fiber.ErrInternalServerError
	}
	balanceSummary := order["balanceSummary"].(map[string]interface{})
	balance := balanceSummary["balance"].(map[string]interface{})
	price := balance["formattedAmount"].(string)
	s := getShippingInfo(order)
	cust := getCustomerInfo(order)
	items := getOrderItems(order)
	details := modalDetails(price, s, cust, items)
	detailsCTX := context.Background()
	var bf bytes.Buffer
	details.Render(detailsCTX, &bf)
	return c.Type("html").Send(bf.Bytes())
}
func getOrderItems(order map[string]interface{}) []item {
	lineItems := order["lineItems"].([]interface{})
	items := make([]item, 0)
	for _, i := range lineItems {
		d := i.(map[string]interface{})
		name := d["productName"].(map[string]interface{})["original"].(string)
		catalogRef := d["catalogReference"].(map[string]interface{})
		image := d["image"].(map[string]interface{})["url"].(string)
		price := d["price"].(map[string]interface{})["formattedAmount"].(string)
		variantsOpt, ok := catalogRef["options"]
		if ok {
			colorVar, ok := variantsOpt.(map[string]interface{})["options"].(map[string]interface{})["Color"]
			if ok {
				name = fmt.Sprintf("%s - %s", name, colorVar)
			}

			sizeVar, ok := variantsOpt.(map[string]interface{})["options"].(map[string]interface{})["Size"]
			if ok {
				name = fmt.Sprintf("%s - %s", name, sizeVar)
			}
		}
		items = append(
			items,
			item{title: name, price: price, imageUrl: image, quantity: d["quantity"].(float64)},
		)
	}
	return items
}
func getCustomerInfo(order map[string]interface{}) customer {
	shippingInfo := order["shippingInfo"].(map[string]interface{})["logistics"].(map[string]interface{})["shippingDestination"].(map[string]interface{})
	contactDetails := shippingInfo["contactDetails"].(map[string]interface{})
	firstName := contactDetails["firstName"].(string)
	lastName := contactDetails["lastName"].(string)
	phoneNumber := contactDetails["phone"].(string)
	c := customer{name: fmt.Sprintf("%s %s", firstName, lastName), phone: phoneNumber}
	return c
}
func getShippingInfo(order map[string]interface{}) shipping {
	shiippingPrice := order["priceSummary"].(map[string]interface{})["shipping"].(map[string]interface{})["formattedAmount"].(string)
	shippingInfo := order["shippingInfo"].(map[string]interface{})["logistics"].(map[string]interface{})["shippingDestination"].(map[string]interface{})
	address := shippingInfo["address"].(map[string]interface{})
	coutry := address["country"].(string)
	city := address["city"].(string)
	zipCode := address["postalCode"].(string)
	addressLine := address["addressLine"].(string)
	s := shipping{
		deliveryPrice:   shiippingPrice,
		deliveryAddress: fmt.Sprintf("%s.\n%s, %s, %s", addressLine, city, zipCode, coutry),
	}
	return s
}
func OrderPlacedEvent(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	webhookData := make(map[string]interface{})
	err := c.BodyParser(&webhookData)
	if err != nil {
		log.Error(err)
		return fiber.ErrBadRequest
	}
	data := webhookData["data"].(map[string]interface{})
	rawDate := data["createdDate"].(string)
	creationDate, err := time.Parse(time.RFC3339, rawDate)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	platformID, _ := data["id"].(string)
	totalPrice, _ := data["priceSummary"].(map[string]interface{})["total"].(map[string]interface{})["value"].(string)
	Price, _ := strconv.ParseFloat(totalPrice, 64)
	firstName, _ := data["contact"].(map[string]interface{})["name"].(map[string]interface{})["first"].(string)
	lastName, _ := data["contact"].(map[string]interface{})["name"].(map[string]interface{})["last"].(string)
	orderNumber, _ := strconv.Atoi(data["orderNumber"].(string))
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	_, err = q.AddOrder(dbCTX, orders.AddOrderParams{
		PlatformID: platformID, CreationDate: creationDate.Unix(), OrderNumber: int64(orderNumber), TotalPrice: Price, CustomerName: fmt.Sprintf("%s %s", lastName, firstName),
	})
	if err != nil {
		log.Error(err)
		return fiber.NewError(fiber.StatusBadRequest, "order already exists on DB")
	}
	return c.SendStatus(fiber.StatusOK)
}
