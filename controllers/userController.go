package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krushev/go-footing/models"
	"github.com/krushev/go-footing/services"
	"github.com/krushev/go-footing/util"
	"net/http"
)

type userController struct {
	us services.UserService
}

func NewUserController(us services.UserService) userController {
	return userController{us}
}

func (uc *userController) Router() *fiber.App {
	app := fiber.New()
	app.Get("/", uc.findActive)
	app.Get("/search", uc.search)
	app.Get("/:id", uc.get)
	app.Post("/", uc.create)
	app.Put("/:id", uc.put)
	app.Delete("/:id", uc.del)

	return app
}

func (uc *userController) create(ctx *fiber.Ctx) error {
	var newUser models.User
	err := ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(util.NewJError(err))
	}
	user, err := uc.us.Create(newUser)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(err))
	}
	return ctx.
		Status(http.StatusCreated).
		JSON(user)
}

func (uc *userController) findActive(ctx *fiber.Ctx) error {
	users, err := uc.us.FindActive()
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(err))
	}
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"data":  users,
			"total": len(*users),
		})
}

func (uc *userController) search(ctx *fiber.Ctx) error {
	pagination := Pagination(ctx)
	users, err := uc.us.Search(ctx.Query("q"), &pagination)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(err))
	}
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"data":  users,
			"total": len(*users),
		})
}

func (uc *userController) get(c *fiber.Ctx) error {
	user, err := uc.us.GetById(c.Params("id"))
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(user)
}

func (uc *userController) put(ctx *fiber.Ctx) error {
	return nil
}

func (uc *userController) del(ctx *fiber.Ctx) error {
	return nil
}
