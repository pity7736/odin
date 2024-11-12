package categoryhandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler/restcategoryhandler"
)

type CategoryHandler interface {
	CreateCommand() categorycommand.CategoryCreatorCommand
	HandleOneResponse(category *category.Category)
	HandleManyResponse(categories []*category.Category)
}

type categoryHandler struct {
	repository repositories.CategoryRepository
	handler    CategoryHandler
}

func New(repository repositories.CategoryRepository) *categoryHandler {
	return &categoryHandler{repository: repository}
}

func (c *categoryHandler) Post(ctx *fiber.Ctx) error {
	c.setHandler(ctx)
	command := c.handler.CreateCommand()
	categoryCreator := categorycreator.New(command, c.repository)
	category, _ := categoryCreator.Create(ctx.Context())
	c.handler.HandleOneResponse(category)
	ctx.Status(http.StatusCreated)
	return nil
}

func (c *categoryHandler) Get(ctx *fiber.Ctx) error {
	c.setHandler(ctx)
	categories := c.repository.GetAll(ctx.Context())
	c.handler.HandleManyResponse(categories)
	return nil
}

func (c *categoryHandler) setHandler(ctx *fiber.Ctx) {
	c.handler = restcategoryhandler.New(ctx)
}
