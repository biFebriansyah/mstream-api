package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func New(*sqlx.DB) *fiber.App {
	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("hello worlds")
	})

	return app
}
