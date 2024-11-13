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

func New(ctx *fiber.Ctx) *restCategoryHandler {
	return &restCategoryHandler{ctx: ctx}
}

func (r *restCategoryHandler) CreateCommand() categorycommand.CategoryCreatorCommand {
	var body categoryrequestbody.CategoryRequestBody
	r.ctx.BodyParser(&body)
	return body.CreateCategoryCreatorCommand()
}

func (r *restCategoryHandler) HandleOneResponse(category *category.Category) {
	r.ctx.JSON(r.getCategoryResponse(category))
}

func (r *restCategoryHandler) HandleManyResponse(categories []*category.Category) {
	result := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, r.getCategoryResponse(category))
	}
	r.ctx.JSON(CategoriesResponse{Categories: result})
}

func (r *restCategoryHandler) getCategoryResponse(category *category.Category) CategoryResponse {
	return CategoryResponse{
		Id:   category.ID(),
		Name: category.Name(),
		Type: category.Type().String(),
		User: category.User().Email(),
	}
}

type CategoryResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	User string `json:"user"`
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}