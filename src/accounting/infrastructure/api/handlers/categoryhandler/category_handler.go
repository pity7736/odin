package categoryhandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type CategoryHandler struct {
	repository repositories.CategoryRepository
}

func New(repository repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repository: repository}
}

func (c *CategoryHandler) Handle(ctx *fiber.Ctx) error {
	var body categoryBody
	ctx.BodyParser(&body)
	command := body.CreateCategoryCreatorCommand()
	categoryCreator := categorycreator.New(command, c.repository)
	categoryCreator.Create(ctx.Context())
	ctx.Status(http.StatusCreated)
	return nil
}
