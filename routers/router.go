package routers

import (
	"biFebriansyah/gostream/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *fiber.App {
	app := fiber.New(config.FiberConfig)
	app.Use(recover.New())

	artisRoute(app, db)

	return app
}
