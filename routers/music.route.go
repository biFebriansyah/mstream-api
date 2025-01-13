package routers

import (
	"biFebriansyah/gostream/handlers"
	"biFebriansyah/gostream/middleware"
	"biFebriansyah/gostream/repositories"
	"biFebriansyah/gostream/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func musicRoute(app *fiber.App, db *sqlx.DB, amqp *utils.AmqpConfig) {
	music := app.Group("/music")

	repos := repositories.NewMusic(db)
	handler := handlers.NewMusicHandler(repos, amqp)

	music.Get("/status", middleware.SockCheck, utils.FFmpegStatus())
	music.Get("/slug/:slug", handler.FetchBySlug)
	music.Get("/:uuid", handler.FetchById)
	music.Get("/", handler.FetchAll)
	music.Post("/", middleware.Upload(), handler.Create)
	music.Put("/", handler.Update)
	music.Delete("/:uuid", handler.Delete)

}
