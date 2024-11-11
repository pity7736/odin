package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"raiseexception.dev/odin/src/shared/infrastructure/repositoryfactory"
)

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(repositoryFactory repositoryfactory.RepositoryFactory) Application {
	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Post("/v1/categories", categoryhandler.New(repositoryFactory.GetCategoryRepository()).Handle)
	return &fibberApplication{app: app}
}

func (a *fibberApplication) Start() error {
	return a.app.Listen(":8000")
}

func (a *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return a.app.Test(request)
}
