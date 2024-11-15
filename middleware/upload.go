package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	FormName []string
	DirName  string
}

var ConfigDefault = Config{
	FormName: []string{"file"},
	DirName:  "./uploads",
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]
	return cfg
}

func Upload(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	return func(c *fiber.Ctx) (err error) {
		// TODO upload multiple
		for _, v := range cfg.FormName {
			file, err := c.FormFile(v)
			if err != nil {
				return err
			}

			fileLocation := fmt.Sprintf("%s/%s", cfg.DirName, file.Filename)
			if err := c.SaveFile(file, fileLocation); err != nil {
				return err
			}

			c.Locals("file", fileLocation)

		}

		return c.Next()
	}
}
