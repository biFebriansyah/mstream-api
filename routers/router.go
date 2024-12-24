package routers

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, amqp *utils.AmqpConfig) *fiber.App {
	app := fiber.New(config.FiberConfig)
	app.Use(cors.New())
	app.Use(recover.New())

	artisRoute(app, db)
	genreRoute(app, db)
	musicRoute(app, db, amqp)

	return app
}
