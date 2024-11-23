package restcategoryhandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
)

type restCategoryHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) categoryhandler.CategoryHandler {
	return &restCategoryHandler{ctx: ctx}
}

func (self *restCategoryHandler) HandleOneResponse(category *categorymodel.Category) {
	self.ctx.JSON(self.getCategoryResponse(category))
}

func (self *restCategoryHandler) HandleManyResponse(categories []*categorymodel.Category) {
	result := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, self.getCategoryResponse(category))
	}
	self.ctx.JSON(CategoriesResponse{Categories: result})
}

func (self *restCategoryHandler) getCategoryResponse(category *categorymodel.Category) CategoryResponse {
	return CategoryResponse{
		Id:     category.ID(),
		Name:   category.Name(),
		Type:   category.Type().String(),
		UserID: category.UserID(),
	}
}

func (self *restCategoryHandler) ContentType() string {
	return fiber.MIMEApplicationJSON
}

type CategoryResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	UserID string `json:"user_id"`
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}
