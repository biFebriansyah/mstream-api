package routers

import (
	"biFebriansyah/gostream/handlers"
	"biFebriansyah/gostream/middleware"
	"biFebriansyah/gostream/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func artisRoute(app *fiber.App, db *sqlx.DB) {
	artis := app.Group("/artis")

	uploadConfig := middleware.Config{FormName: []string{"picture"}}
	repos := repositories.NewArtis(db)
	handler := handlers.NewartisHandler(repos)

	artis.Get("/slug/:slug", handler.FetchBySlug)
	artis.Get("/:uuid", handler.FetchById)
	artis.Get("/", handler.FetchAll)
	artis.Post("/", middleware.Upload(uploadConfig), handler.Create)
	artis.Put("/", handler.Update)
	artis.Delete("/:uuid", handler.Delete)

}
