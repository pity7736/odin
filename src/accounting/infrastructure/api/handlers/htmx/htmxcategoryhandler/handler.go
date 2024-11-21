package htmxcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"strconv"
)

type htmxCategoryHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) *htmxCategoryHandler {
	return &htmxCategoryHandler{ctx: ctx}
}

func (h *htmxCategoryHandler) HandleOneResponse(category *category.Category) {
	isFirst, _ := strconv.ParseBool(h.ctx.FormValue("first", "false"))
	if isFirst {
		h.ctx.Set("HX-Refresh", "true")
	}
	h.ctx.Render("category", category, "")
}

func (h *htmxCategoryHandler) HandleManyResponse(categories []*category.Category) {
	h.ctx.Render("categories", Data{Categories: categories})
}

func (h *htmxCategoryHandler) ContentType() string {
	return fiber.MIMETextHTMLCharsetUTF8
}

type Data struct {
	Categories []*category.Category
	Category   *category.Category
}
