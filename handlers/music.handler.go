package handlers

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"biFebriansyah/gostream/utils"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberutil "github.com/gofiber/fiber/v2/utils"
)

type musicHandler struct {
	repo *repositories.MusicRepo
	amqp *utils.AmqpConfig
}

func NewMusicHandler(repo *repositories.MusicRepo, amqp *utils.AmqpConfig) *musicHandler {
	return &musicHandler{repo, amqp}
}

func (music *musicHandler) Create(ctx *fiber.Ctx) error {
	data := new(models.MusicData)
	clean := utils.Cleaning()
	upload := utils.NewGIO()

	if err := ctx.BodyParser(data); err != nil {
		log.Println(err)
		return fiber.ErrBadGateway
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println(err)
		return fiber.ErrBadGateway
	}

	if artisData, oke := form.Value["music_artis"]; oke {
		if err := json.Unmarshal([]byte(artisData[0]), &data.MusicArtis); err != nil {
			return fiber.ErrBadRequest
		}
	}

	if genreData, oke := form.Value["music_genre"]; oke {
		if err := json.Unmarshal([]byte(genreData[0]), &data.MusicGenre); err != nil {
			return fiber.ErrBadRequest
		}
	}

	if file := ctx.Locals("image").(string); file != "" {
		if url, err := upload.UploadData(file); err == nil {
			data.Cover = url
			clean.Add(file)
		}
	} else {
		return fiber.ErrBadGateway
	}

	data.Slug = utils.Slug(data.Title)
	result, err := music.repo.InsertData(data)
	if err != nil {
		return err
	}

	if file := ctx.Locals("file").(string); file != "" {
		message := map[string]string{"uuid": result, "location": file}
		if err := music.amqp.NewPublisher("ffmpeg", message); err != nil {
			return fiber.ErrBadGateway
		}
	} else {
		return fiber.ErrBadGateway
	}

	go clean.Run()
	// return ctx.JSON(utils.Respone(result))
	return ctx.JSON(utils.Respone(fmt.Sprintf("%s data created", result)))
}

func (music *musicHandler) Update(ctx *fiber.Ctx) error {
	data := new(models.Music)

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	if data.Title != "" {
		data.Slug = utils.Slug(data.Title)
	}

	result, err := music.repo.UpdateData(data)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data updated", result)))
}

func (music *musicHandler) Delete(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := music.repo.DeleteData(uid)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data delete", result)))
}

func (music *musicHandler) FetchById(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := music.repo.GetDataById(uid)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (music *musicHandler) FetchBySlug(ctx *fiber.Ctx) error {
	slug := fiberutil.CopyString(ctx.Params("slug"))
	result, err := music.repo.GetDataBySlug(slug)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (music *musicHandler) FetchAll(ctx *fiber.Ctx) error {
	paginate := &config.Pagination{Page: 1, Limit: 10}
	paginate.Title = ctx.Query("title")

	if page, err := strconv.Atoi(ctx.Query("page", "1")); err == nil {
		paginate.Page = int32(page)
	}
	if limit, err := strconv.Atoi(ctx.Query("limit", "10")); err == nil {
		paginate.Limit = int32(limit)
	}

	result, err := music.repo.GetAllData(paginate)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}
