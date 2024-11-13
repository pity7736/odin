package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"raiseexception.dev/odin/src/shared/infrastructure/repositoryfactory"
)

const categoriesPath = "/categories"

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(repositoryFactory repositoryfactory.RepositoryFactory) Application {
	engine := html.New("./src/shared/infrastruture/templates", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base",
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	v1 := app.Group("/v1")
	v1.Post(categoriesPath, func(c *fiber.Ctx) error {
		return categoryhandler.New(repositoryFactory.GetCategoryRepository()).Create(c)
	})
	v1.Get(categoriesPath, func(c *fiber.Ctx) error {
		return categoryhandler.New(repositoryFactory.GetCategoryRepository()).GetAll(c)
	})
	return &fibberApplication{app: app}
}

func (a *fibberApplication) Start() error {
	return a.app.Listen(":8000")
}

func (a *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return a.app.Test(request)
}
