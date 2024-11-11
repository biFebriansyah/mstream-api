package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type response struct {
	Status      int         `json:"status"`
	IsError     bool        `json:"isError"`
	Description interface{} `json:"description"`
}

type ResultWarp struct {
	Data any
	Meta any
}

type Meta struct {
	Next  int32 `json:"next"`
	Prev  int32 `json:"prev"`
	Total int32 `json:"total"`
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
