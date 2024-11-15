package htmxcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler/categoryrequestbody"
	"strconv"
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
	h.ctx.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	isFirst, _ := strconv.ParseBool(h.ctx.FormValue("first", "false"))
	if isFirst {
		h.ctx.Set("HX-Refresh", "true")
	}
	h.ctx.Render("category", category, "")
}

func (h *htmxCategoryHandler) HandleManyResponse(categories []*category.Category) {
	h.ctx.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	h.ctx.Render("categories", Data{Categories: categories})
}

type Data struct {
	Categories []*category.Category
	Category   *category.Category
}
