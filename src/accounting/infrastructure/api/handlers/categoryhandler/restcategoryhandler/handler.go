package restcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler/categoryrequestbody"
)

type restCategoryHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) restCategoryHandler {
	return restCategoryHandler{ctx: ctx}
}

func (r restCategoryHandler) CreateCommand() categorycommand.CategoryCreatorCommand {
	var body categoryrequestbody.CategoryRequestBody
	r.ctx.BodyParser(&body)
	return body.CreateCategoryCreatorCommand()
}

func (r restCategoryHandler) HandleResponse(category *category.Category) {
	response := categoryResponse{
		Id:   category.ID(),
		Name: category.Name(),
		Type: category.Type().String(),
		User: category.User().Email(),
	}
	r.ctx.JSON(response)
}

type categoryResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	User string `json:"user"`
}
