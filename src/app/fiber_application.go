package app

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/htmx/htmxcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/accountsrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/htmx/htmxloginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/rest/restloginhandler"
)

const categoriesPath = "/categories"

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(accountingRepositoryFactory accountingrepositoryfactory.RepositoryFactory,
	accountsRepositoryFactory accountsrepositoryfactory.AccountsRepositoryFactory,
) Application {
	engine := html.New(
		"/Users/julian.cortes/development/odin/src/shared/infrastructure/templates",
		".gohtml",
	)
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base",
	})
	app.Use(logger.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Use(func(c *fiber.Ctx) error {
		cookie := c.Cookies("__Secure-odin-session")
		if cookie != "" {
			session, _ := accountsRepositoryFactory.GetSessionRepository().Get(c.Context(), cookie)
			c.Locals("userID", session.UserID())
		}
		return c.Next()
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	api := app.Group("/api")
	apiV1 := api.Group("/v1")
	apiV1.Use(func(c *fiber.Ctx) error {
		authHeader := strings.Split(c.Get("Authorization", ""), " ")
		if len(authHeader) == 2 {
			session, _ := accountsRepositoryFactory.GetSessionRepository().Get(c.Context(), authHeader[1])
			c.Locals("userID", session.UserID())
		}
		return c.Next()
	})
	apiV1.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		if ctx.Locals("userID") != nil {
			// TODO: handle error. design a way to handle app errors
			categoryhandler.New(
				accountingRepositoryFactory.GetCategoryRepository(),
				restcategoryhandler.New(ctx),
			).Create(ctx)
			return nil
		} else {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}
	})
	apiV1.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		if ctx.Locals("userID") != nil {
			return categoryhandler.New(
				accountingRepositoryFactory.GetCategoryRepository(),
				restcategoryhandler.New(ctx),
			).GetAll(ctx)
		} else {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}
	})
	app.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).Create(ctx)
		return nil
	})
	app.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).GetAll(ctx)
		return nil
	})
	apiV1.Post("/auth/login", func(ctx *fiber.Ctx) error {
		return loginhandler.New(
			accountsRepositoryFactory,
			restloginhandler.New(ctx),
		).Login(ctx)
	})
	app.Get("/auth/login", func(ctx *fiber.Ctx) error {
		ctx.Render("login", htmxloginhandler.RequestError{Error: ""})
		return nil
	})
	app.Post("/auth/login", func(ctx *fiber.Ctx) error {
		return loginhandler.New(
			accountsRepositoryFactory,
			htmxloginhandler.New(ctx),
		).Login(ctx)
	})
	return &fibberApplication{app: app}
}

func (self *fibberApplication) Start() error {
	return self.app.Listen(":8000")
}

func (self *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return self.app.Test(request)
}
