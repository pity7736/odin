package htmxcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler/categoryrequestbody"
)

type htmxCategoryHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) *htmxCategoryHandler {
	return &htmxCategoryHandler{ctx: ctx}
}

func (h *htmxCategoryHandler) CreateCommand() categorycommand.CategoryCreatorCommand {
	var body categoryrequestbody.CategoryRequestBody
	h.ctx.BodyParser(&body)
	return body.CreateCategoryCreatorCommand()
}

func (h *htmxCategoryHandler) HandleOneResponse(category *category.Category) {

}

func (h *htmxCategoryHandler) HandleManyResponse(categories []*category.Category) {
	h.ctx.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	h.ctx.Render("categories", categories, "")
}
