package categoryhandler

import (
	"net/http"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler/categoryrequestbody"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/htmx/htmxcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
)

type CategoryHandler interface {
	HandleOneResponse(category *category.Category)
	HandleManyResponse(categories []*category.Category)
	ContentType() string
}

type categoryHandler struct {
	repository repositories.CategoryRepository
	handler    CategoryHandler
}

func New(repository repositories.CategoryRepository) *categoryHandler {
	return &categoryHandler{repository: repository}
}

func (c *categoryHandler) Create(ctx *fiber.Ctx) error {
	c.setHandler(ctx)
	ctx.Set("Content-Type", c.handler.ContentType())
	command, errCmd := c.createCommand(ctx)
	if errCmd != nil {
		ctx.Status(http.StatusBadRequest)
		return errCmd
	}
	categoryCreator := categorycreator.New(*command, c.repository)
	category, _ := categoryCreator.Create(ctx.Context())
	c.handler.HandleOneResponse(category)
	ctx.Status(http.StatusCreated)
	return nil
}

func (c *categoryHandler) createCommand(ctx *fiber.Ctx) (*categorycommand.CategoryCreatorCommand, error) {
	var body categoryrequestbody.CategoryRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return nil, err
	}
	return body.CreateCategoryCreatorCommand()
}

func (c *categoryHandler) GetAll(ctx *fiber.Ctx) error {
	c.setHandler(ctx)
	categories := c.repository.GetAll(ctx.Context())
	c.handler.HandleManyResponse(categories)
	return nil
}

func (c *categoryHandler) setHandler(ctx *fiber.Ctx) {
	if ctx.Get("content-type") == fiber.MIMEApplicationJSON {
		c.handler = restcategoryhandler.New(ctx)
	} else {
		c.handler = htmxcategoryhandler.New(ctx)
	}
}
