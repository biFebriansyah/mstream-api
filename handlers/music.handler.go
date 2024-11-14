package handlers

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"biFebriansyah/gostream/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberutil "github.com/gofiber/fiber/v2/utils"
)

type musicHandler struct {
	repo *repositories.MusicRepo
}

func NewMusicHandler(repo *repositories.MusicRepo) *musicHandler {
	return &musicHandler{repo}
}

func (music *musicHandler) Create(ctx *fiber.Ctx) error {
	data := new(models.Music)

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	data.Slug = utils.Slug(data.Title)
	result, err := music.repo.InsertData(data)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data created", result)))
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
		log.Println(err)
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}
