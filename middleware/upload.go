package middleware

import (
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Config struct {
	FormName []string
	dirName  string
}

var ConfigDefault = Config{
	FormName: []string{"file", "image"},
	dirName:  "./uploads",
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]
	cfg.dirName = ConfigDefault.dirName
	return cfg
}

func Upload(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	return func(c *fiber.Ctx) (err error) {
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, v := range cfg.FormName {
			uidds := uuid.New().String()[:8]
			files := form.File[v]
			for _, file := range files {
				fileExt := filepath.Ext(file.Filename)
				fileName := fmt.Sprintf("%s%s", uidds, fileExt)
				fileLocation := fmt.Sprintf("%s/%s", cfg.dirName, fileName)
				if err := c.SaveFile(file, fileLocation); err != nil {
					return err
				}

				c.Locals(v, fileLocation)
			}

		}

		return c.Next()
	}
}
