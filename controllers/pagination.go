package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krushev/go-footing/models"
	"go.uber.org/zap"
	"strconv"
)

func Pagination(ctx *fiber.Ctx) models.Pagination {
	result := ctx.Params("page")
	zap.S().Infof(result)

	size, _ := strconv.Atoi(ctx.Params("size", "5"))
	page, _ := strconv.Atoi(ctx.Params("page", "1"))
	sort := ctx.Params("sort", `created_at DESC`)

	return models.Pagination{
		Size: size,
		Page: page,
		Sort: sort,
	}
}
