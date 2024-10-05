package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/braxes-backend/app/middleware"
	"github.com/braxes-backend/database"
	"github.com/braxes-backend/database/orders"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// var queries *orders.Queries

type reqOrder struct {
	PlatformID int64
}
type respOrder struct {
	Id           int64 `json:"id"`
	PlatformID   int64 `json:"platformID"`
	CreationDate int64 `json:"creationDate"`
	IsProcessed  int64 `json:"isProcessed"`
}
type _ orders.Order

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
	handler.Patch("/unprocess/:id", UnProcessOrder)
	handler.Patch("/process/:id", ProcessOrder)
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
	ordersResp := make([]respOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			respOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
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
	orders, err := q.GetUnProcessedOrders(dbCTX)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]respOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			respOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
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
func ListProcessed(c *fiber.Ctx) error {
	q, dbCTX := middleware.GetQueryCTX(c)
	orders, err := q.GetProcessedOrders(dbCTX)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	ordersResp := make([]respOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			respOrder{
				Id:           o.ID,
				PlatformID:   o.PlatformID,
				CreationDate: o.CreationDate,
				IsProcessed:  o.IsProcessed,
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
		respOrder{
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
	ordersResp := make([]respOrder, 0)
	for _, o := range orders {
		ordersResp = append(
			ordersResp,
			respOrder{
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Errorf("Cannot parse id from path params,%e\n", err)
		return fiber.ErrInternalServerError
	}
	q, dbCTX := middleware.GetQueryCTX(c)
	processed, err := q.ProcessOrder(
		dbCTX,
		orders.ProcessOrderParams{
			ProcessedDate: sql.NullInt64{Valid: true, Int64: time.Now().Unix()}, ID: int64(id),
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Errorf("id cannot be parsed to int, %e\n", err)
		return fiber.ErrInternalServerError
	}
	q, dbCTX := middleware.GetQueryCTX(c)
	processed, err := q.UnProcessOrder(
		dbCTX,
		int64(id),
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Error("cannot parse id param")
		return fiber.ErrBadRequest
	}
	q, dbCTX := middleware.GetQueryCTX(c)
	order, err := q.GetOrder(dbCTX, int64(id))
	if err != nil {
		log.Error(err)
		return fiber.ErrNotFound
	}
	resp := respOrder{
		Id:           order.ID,
		PlatformID:   order.PlatformID,
		IsProcessed:  order.IsProcessed,
		CreationDate: order.CreationDate,
	}
	json, err := json.Marshal(resp)
	return c.Send(json)
}
func OrderPlacedEvent(c *fiber.Ctx) error { return nil }
