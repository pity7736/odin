package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxcreateaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxgetaccountshandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/restcreateaccounthandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	"raiseexception.dev/odin/src/shared/infrastructure/api"

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
const accountPath = "/accounts"

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
		Views:        engine,
		ViewsLayout:  "base",
		ErrorHandler: errorHandler,
	})
	app.Use(logger.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Use(func(c *fiber.Ctx) error {
		cookie := c.Cookies("__Secure-odin-session")
		c.Locals(requestcontext.Key, requestcontext.NewAnonymous())
		if cookie != "" {
			session, _ := accountsRepositoryFactory.GetSessionRepository().Get(c.Context(), cookie)
			if session != nil {
				c.Locals("userID", session.UserID())
				requestContext, err := requestcontext.New(session.UserID())
				if err != nil {
					return err
				}
				c.Locals(requestcontext.Key, requestContext)
			}
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
		if ctx.Locals("userID") != nil {
			categoryhandler.New(
				accountingRepositoryFactory.GetCategoryRepository(),
				htmxcategoryhandler.New(ctx),
			).Create(ctx)
			return nil
		} else {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}
	})
	app.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		if ctx.Locals("userID") != nil {
			categoryhandler.New(
				accountingRepositoryFactory.GetCategoryRepository(),
				htmxcategoryhandler.New(ctx),
			).GetAll(ctx)
			return nil
		} else {
			ctx.Status(http.StatusUnauthorized)
			return nil
		}
	})
	apiV1.Post("/auth/login", func(ctx *fiber.Ctx) error {
		return loginhandler.New(
			accountsRepositoryFactory,
			restloginhandler.New(ctx),
		).Login(ctx)
	})
	apiV1.Post(accountPath, func(ctx *fiber.Ctx) error {
		return restcreateaccounthandler.New(accountingRepositoryFactory.GetAccountRepository()).Handle(ctx)
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
	app.Post(accountPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxcreateaccounthandler.New(accountingRepositoryFactory.GetAccountRepository()))
	})
	app.Get(accountPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxgetaccountshandler.New(accountingRepositoryFactory.GetAccountRepository()))
	})
	return &fibberApplication{app: app}
}

func (self *fibberApplication) Start() error {
	return self.app.Listen(":8000")
}

func (self *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return self.app.Test(request, -1)
}

func loginRequired(ctx *fiber.Ctx, handler api.Handler) error {
	requestContext := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	if requestContext.IsAuthenticated() {
		return handler.Handle(ctx)
	}
	ctx.Status(http.StatusUnauthorized)
	ctx.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
	if ctx.Get("Content-Type", "") == fiber.MIMEApplicationJSON {
		ctx.Set("Content-Type", fiber.MIMEApplicationJSON)
	}
	return nil
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	var odinError *odinerrors.Error
	ok := errors.As(err, &odinError)
	code := http.StatusInternalServerError
	if ok {
		switch odinError.Tag() {
		case odinerrors.DOMAIN:
			code = http.StatusBadRequest
		default:
		}
	}
	ctx.Status(code)
	return nil
}
