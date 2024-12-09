package handlers

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"biFebriansyah/gostream/repositories"
	"biFebriansyah/gostream/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberutil "github.com/gofiber/fiber/v2/utils"
)

type artisHandler struct {
	repo *repositories.ArtiRepo
}

func NewartisHandler(repo *repositories.ArtiRepo) *artisHandler {
	return &artisHandler{repo}
}

func (artis *artisHandler) Create(ctx *fiber.Ctx) error {
	data := new(models.Artis)
	clean := utils.Cleaning()
	upload := utils.NewGIO()

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	if file := ctx.Locals("picture").(string); file != "" {
		if url, err := upload.UploadData(file); err == nil {
			data.Picture = url
			clean.Add(file)
		}
	}

	data.Slug = utils.Slug(data.FirstName + data.LastName)
	result, err := artis.repo.InsertData(data)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data created", result)))
}

func (artis *artisHandler) Update(ctx *fiber.Ctx) error {
	data := new(models.Artis)

	if err := ctx.BodyParser(data); err != nil {
		return fiber.ErrBadGateway
	}

	if data.FirstName != "" || data.LastName != "" {
		data.Slug = utils.Slug(data.FirstName + data.LastName)
	}

	result, err := artis.repo.UpdateData(data)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data updated", result)))
}

func (artis *artisHandler) Delete(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := artis.repo.DeleteData(uid)
	if err != nil {
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(fmt.Sprintf("%d data delete", result)))
}

func (artis *artisHandler) FetchById(ctx *fiber.Ctx) error {
	uid := fiberutil.CopyString(ctx.Params("uuid"))
	result, err := artis.repo.GetDataById(uid)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (artis *artisHandler) FetchBySlug(ctx *fiber.Ctx) error {
	slug := fiberutil.CopyString(ctx.Params("slug"))
	result, err := artis.repo.GetDataBySlug(slug)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}

func (artis *artisHandler) FetchAll(ctx *fiber.Ctx) error {
	var paginate = &config.Pagination{Page: 1, Limit: 10}
	paginate.Name = ctx.Query("name")

	if page, err := strconv.Atoi(ctx.Query("page", "1")); err != nil {
		paginate.Page = int32(page)
	}
	if limit, err := strconv.Atoi(ctx.Query("limit", "10")); err != nil {
		paginate.Limit = int32(limit)
	}

	result, err := artis.repo.GetAllData(paginate)
	if err != nil {
		if err.Error() == config.NotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrBadGateway
	}

	return ctx.JSON(utils.Respone(result))
}
