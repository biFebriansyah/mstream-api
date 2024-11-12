package routers

import (
	"biFebriansyah/gostream/handlers"
	"biFebriansyah/gostream/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func genreRoute(app *fiber.App, db *sqlx.DB) {
	genre := app.Group("/genre")

	repos := repositories.NewGenre(db)
	handler := handlers.NewGenreHandler(repos)

	genre.Get("/slug/:slug", handler.FetchBySlug)
	genre.Get("/:uuid", handler.FetchById)
	genre.Get("/", handler.FetchAll)
	genre.Post("/", handler.Create)
	genre.Put("/", handler.Update)
	genre.Delete("/:uuid", handler.Delete)

}
