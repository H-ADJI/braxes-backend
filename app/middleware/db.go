package middleware

import (
	"context"

	"github.com/braxes-backend/database/orders"
	"github.com/gofiber/fiber/v2"
)

func InjectDB(q *orders.Queries, dbCTX context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("queries", q)
		c.Locals("dbCTX", dbCTX)
		return c.Next()
	}
}
func GetQueryCTX(c *fiber.Ctx) (*orders.Queries, context.Context) {
	q := GetLocal[*orders.Queries](c, "queries")
	ctx := GetLocal[context.Context](c, "dbCTX")
	return q, ctx

}
