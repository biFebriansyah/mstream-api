package utils

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewServer(route *fiber.App) *fiber.App {
	var addrs string = "0.0.0.0:8080"
	if pr := os.Getenv("PORT"); pr != "" {
		addrs = ":" + pr
	}

	fiberServer := fiber.New()
	fiberServer.Mount("/v1", route)

	go func() {
		if err := fiberServer.Listen(addrs); err != nil {
			log.Fatalf("Could not listen on %s: %v\n", addrs, err)
		}
	}()

	return fiberServer
}
