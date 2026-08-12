package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxcreateaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxgetaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxgetaccountshandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/restcreateaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/categoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/htmx/htmxcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/incomehandler/htmxcreateincomehandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"

	"raiseexception.dev/odin/src/accounts/application/passwordhasher"
	accountsrepos "raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/htmx/htmxloginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/htmx/htmxlogouthandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/logouthandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/rest/restloginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/rest/restlogouthandler"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

const categoriesPath = "/categories"
const accountPath = "/accounts"

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(
	accountingRepositoryFactory accountingrepositoryfactory.RepositoryFactory,
	sessionRepository accountsrepos.SessionRepository,
	userRepository accountsrepos.UserRepository,
	passwordHasher passwordhasher.PasswordHasher,
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
	app.Use(cookieMiddleware(sessionRepository))
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	app.Get("/auth/login", func(ctx *fiber.Ctx) error {
		next := ctx.Query("next", "/")
		return ctx.Render("login", htmxloginhandler.LoginData{Error: "", Next: next})
	})
	app.Post("/auth/login", func(ctx *fiber.Ctx) error {
		return loginhandler.New(userRepository, sessionRepository, passwordHasher, htmxloginhandler.New(ctx)).Login(ctx)
	})
	app.Post("/auth/logout", func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, logouthandler.New(sessionRepository, htmxlogouthandler.New(ctx)).Logout)
	})
	app.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).Create)
	})
	app.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			htmxcategoryhandler.New(ctx),
		).GetAll)
	})
	app.Get("/accounts/:accountID", func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxgetaccounthandler.New(accountingRepositoryFactory.GetAccountRepository()).Handle)
	})
	app.Post("/accounts/:accountID/incomes", func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxcreateincomehandler.New(accountingRepositoryFactory).Handle)
	})
	app.Post(accountPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxcreateaccounthandler.New(accountingRepositoryFactory.GetAccountRepository()).Handle)
	})
	app.Get(accountPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, htmxgetaccountshandler.New(accountingRepositoryFactory.GetAccountRepository()).Handle)
	})
	apiV1 := app.Group("/api/v1")
	apiV1.Use(bearerMiddleware(sessionRepository))
	apiV1.Post("/auth/login", func(ctx *fiber.Ctx) error {
		return loginhandler.New(userRepository, sessionRepository, passwordHasher, restloginhandler.New(ctx)).Login(ctx)
	})
	apiV1.Delete("/auth/logout", func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, logouthandler.New(sessionRepository, restlogouthandler.New(ctx)).Logout)
	})
	apiV1.Post(categoriesPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			restcategoryhandler.New(ctx),
		).Create)
	})
	apiV1.Get(categoriesPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, categoryhandler.New(
			accountingRepositoryFactory.GetCategoryRepository(),
			restcategoryhandler.New(ctx),
		).GetAll)
	})
	apiV1.Post(accountPath, func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, restcreateaccounthandler.New(accountingRepositoryFactory.GetAccountRepository()).Handle)
	})
	return &fibberApplication{app: app}
}

func cookieMiddleware(sessionRepository accountsrepos.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(requestcontext.Key, requestcontext.NewAnonymous())
		cookie := c.Cookies("__Secure-odin-session")
		if cookie == "" {
			return c.Next()
		}
		session, err := sessionRepository.Get(c.Context(), cookie)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		if session == nil {
			return c.Next()
		}
		requestCtx, err := requestcontext.New(session.UserID())
		if err != nil {
			return err
		}
		c.Locals(requestcontext.Key, requestCtx)
		c.Locals("userID", session.UserID())
		session.Extend(sessionmodel.DefaultTTL)
		_ = sessionRepository.Save(c.Context(), session)
		return c.Next()
	}
}

func bearerMiddleware(sessionRepository accountsrepos.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Next()
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		session, err := sessionRepository.Get(c.Context(), token)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		if session == nil {
			return c.Next()
		}
		requestCtx, err := requestcontext.New(session.UserID())
		if err != nil {
			return err
		}
		c.Locals(requestcontext.Key, requestCtx)
		c.Locals("userID", session.UserID())
		session.Extend(sessionmodel.DefaultTTL)
		_ = sessionRepository.Save(c.Context(), session)
		return c.Next()
	}
}

func loginRequired(ctx *fiber.Ctx, handlerFn func(*fiber.Ctx) error) error {
	requestCtx := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	if requestCtx.IsAuthenticated() {
		return handlerFn(ctx)
	}
	if strings.HasPrefix(ctx.Path(), "/api/") {
		return ctx.SendStatus(http.StatusUnauthorized)
	}
	ctx.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
	return ctx.Redirect("/auth/login?next=" + ctx.Path())
}

func (self *fibberApplication) Start() error {
	return self.app.Listen(":8000")
}

func (self *fibberApplication) Test(request *http.Request) (*http.Response, error) {
	return self.app.Test(request, -1)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	var odinError *odinerrors.Error
	code := http.StatusInternalServerError
	ok := errors.As(err, &odinError)
	if ok {
		switch odinError.Tag() {
		case odinerrors.DOMAIN:
			code = http.StatusBadRequest
		case odinerrors.NotFound:
			code = http.StatusNotFound
		default:
		}
	}
	ctx.Status(code)
	return nil
}
