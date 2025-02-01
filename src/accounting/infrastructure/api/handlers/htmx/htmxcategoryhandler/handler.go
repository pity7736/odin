package htmxcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"strconv"
)

type htmxCategoryHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) categoryhandler.CategoryHandler {
	return &htmxCategoryHandler{ctx: ctx}
}

func (self *htmxCategoryHandler) HandleOneResponse(category *categorymodel.Category) {
	isFirst, _ := strconv.ParseBool(self.ctx.FormValue("first", "false"))
	if isFirst {
		self.ctx.Set("HX-Refresh", "true")
	}
	self.ctx.Render("category", category, "")
}

func (self *htmxCategoryHandler) HandleManyResponse(categories []*categorymodel.Category) {
	self.ctx.Render("categories", Data{Categories: categories})
}

func (self *htmxCategoryHandler) ContentType() string {
	return fiber.MIMETextHTMLCharsetUTF8
}

type Data struct {
	Categories []*categorymodel.Category
}
