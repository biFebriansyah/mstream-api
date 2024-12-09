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

type genreHandler struct {
	repo *repositories.GenreRepo
}

func NewGenreHandler(repo *repositories.GenreRepo) *genreHandler {
	return &genreHandler{repo}
}

func (genre *genreHandler) Create(ctx *fiber.Ctx) error {
	data := new(models.Genre)

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	data.Slug = utils.Slug(data.Name)
	result, err := genre.repo.InsertData(data)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data created", result)))
}

func (genre *genreHandler) Update(ctx *fiber.Ctx) error {
	data := new(models.Genre)

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	if data.Name != "" {
		data.Slug = utils.Slug(data.Name)
	}

	result, err := genre.repo.UpdateData(data)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data updated", result)))
}

func (genre *genreHandler) Delete(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := genre.repo.DeleteData(uid)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data delete", result)))
}

func (genre *genreHandler) FetchById(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := genre.repo.GetDataById(uid)
	if err != nil {
		log.Println(err)
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (genre *genreHandler) FetchBySlug(ctx *fiber.Ctx) error {
	slug := fiberutil.CopyString(ctx.Params("slug"))
	result, err := genre.repo.GetDataBySlug(slug)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (genre *genreHandler) FetchAll(ctx *fiber.Ctx) error {
	paginate := &config.Pagination{Page: 1, Limit: 10}
	paginate.Name = ctx.Query("name")

	if page, err := strconv.Atoi(ctx.Query("page", "1")); err == nil {
		paginate.Page = int32(page)
	}
	if limit, err := strconv.Atoi(ctx.Query("limit", "10")); err == nil {
		paginate.Limit = int32(limit)
	}

	result, err := genre.repo.GetAllData(paginate)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		log.Println(err)
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}
