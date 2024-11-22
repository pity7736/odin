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
	HandleOneResponse(category *categorymodel.Category)
	HandleManyResponse(categories []*categorymodel.Category)
	ContentType() string
}

type categoryHandler struct {
	repository repositories.CategoryRepository
	handler    CategoryHandler
}

func New(repository repositories.CategoryRepository) *categoryHandler {
	return &categoryHandler{repository: repository}
}

func (self *categoryHandler) Create(ctx *fiber.Ctx) error {
	self.setHandler(ctx)
	ctx.Set("Content-Type", self.handler.ContentType())
	command, errCmd := self.createCommand(ctx)
	if errCmd != nil {
		ctx.Status(http.StatusBadRequest)
		return errCmd
	}
	categoryCreator := categorycreator.New(*command, self.repository)
	category, _ := categoryCreator.Create(ctx.Context())
	self.handler.HandleOneResponse(category)
	ctx.Status(http.StatusCreated)
	return nil
}

func (self *categoryHandler) createCommand(ctx *fiber.Ctx) (*categorycommand.CategoryCreatorCommand, error) {
	var body categoryrequestbody.CategoryRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return nil, err
	}
	return body.CreateCategoryCreatorCommand(ctx.Locals("userID").(string))
}

func (self *categoryHandler) GetAll(ctx *fiber.Ctx) error {
	self.setHandler(ctx)
	categories := self.repository.GetAll(ctx.Context())
	self.handler.HandleManyResponse(categories)
	return nil
}

func (self *categoryHandler) setHandler(ctx *fiber.Ctx) {
	if ctx.Get("content-type") == fiber.MIMEApplicationJSON {
		self.handler = restcategoryhandler.New(ctx)
	} else {
		self.handler = htmxcategoryhandler.New(ctx)
	}
}
