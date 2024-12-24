package config

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type response struct {
	Status      int         `json:"status"`
	IsError     bool        `json:"isError"`
	Description interface{} `json:"description"`
}

const (
	NotFound   string = "NotFound"
	BadData    string = "BadRequest"
	UnAutorize string = "UnAutorize"
)

var FiberConfig fiber.Config = fiber.Config{
	AppName:      "stream-api",
	Prefork:      true,
	ServerHeader: "Anonymouse",
	ErrorHandler: errorHandler,
	BodyLimit:    20 * 1024 * 1024,
}

var FiberCors cors.Config = cors.Config{
	AllowOrigins: "*",
	AllowHeaders: "*",
	AllowMethods: strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodHead,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodPatch,
	}, ","),
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	err = ctx.Status(code).JSON(&response{
		Status:      code,
		IsError:     true,
		Description: err.Error(),
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
