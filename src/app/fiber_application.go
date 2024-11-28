package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"net/http"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/htmx/htmxcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
)

const categoriesPath = "/categories"

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(repositoryFactory accountingrepositoryfactory.RepositoryFactory) Application {
	engine := html.New(
		"/Users/julian.cortes/development/odin/src/shared/infrastructure/templates",
		".gohtml",
	)
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base",
	})
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", c.Get("Authorization"))
		return c.Next()
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		categoryhandler.New(
			repositoryFactory.GetCategoryRepository(),
			restcategoryhandler.New(ctx),
		).Create(ctx)
		return nil
	})
	v1.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		return categoryhandler.New(
			repositoryFactory.GetCategoryRepository(),
			restcategoryhandler.New(ctx),
		).GetAll(ctx)
	})
	app.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		categoryhandler.New(
			repositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).Create(ctx)
		return nil
	})
	app.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		categoryhandler.New(
			repositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).GetAll(ctx)
		return nil
	})
	return &fibberApplication{app: app}
}

func (self *fibberApplication) Start() error {
	return self.app.Listen(":8000")
}

func (self *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return self.app.Test(request)
}
