package categoryhandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type categoryResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	User string `json:"user"`
}

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
	category, _ := categoryCreator.Create(ctx.Context())
	response := categoryResponse{
		Id:   category.ID(),
		Name: category.Name(),
		Type: category.Type().String(),
		User: category.User().Email(),
	}
	ctx.JSON(response)
	ctx.Status(http.StatusCreated)
	return nil
}
