package routers

import (
	"biFebriansyah/gostream/handlers"
	"biFebriansyah/gostream/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func artisRoute(app *fiber.App, db *sqlx.DB) {
	artis := app.Group("/artis")

	repos := repositories.NewArtis(db)
	handler := handlers.NewartisHandler(repos)

	artis.Get("/slug/:slug", handler.FetchBySlug)
	artis.Get("/:uuid", handler.FetchById)
	artis.Get("/", handler.FetchAll)
	artis.Post("/", handler.Create)
	artis.Put("/", handler.Update)
	artis.Delete("/:uuid", handler.Delete)

}
