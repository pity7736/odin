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
	HandleResponse(category *category.Category)
}

type categoryHandler struct {
	repository repositories.CategoryRepository
	handler    CategoryHandler
}

func New(repository repositories.CategoryRepository) *categoryHandler {
	return &categoryHandler{repository: repository}
}

func (c *categoryHandler) Handle(ctx *fiber.Ctx) error {
	c.setHandler(ctx)
	command := c.handler.CreateCommand()
	categoryCreator := categorycreator.New(command, c.repository)
	category, _ := categoryCreator.Create(ctx.Context())
	c.handler.HandleResponse(category)
	ctx.Status(http.StatusCreated)
	return nil
}

func (c *categoryHandler) setHandler(ctx *fiber.Ctx) {
	c.handler = restcategoryhandler.New(ctx)
}
