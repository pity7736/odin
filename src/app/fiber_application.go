package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionvalidator"
	accountsrepos "raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/logouthandler"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/registerhandler"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
)

type fibberApplication struct {
	app *fiber.App
}

func NewFiberApplication(
	sessionRepository accountsrepos.SessionRepository,
	userRepository accountsrepos.UserRepository,
	authHasher authhasher.AuthHasher,
) Application {

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})
	app.Use(logger.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	login := loginhandler.New(userRepository, sessionRepository, authHasher)
	logout := logouthandler.New(sessionRepository)
	register := registerhandler.New(userRepository, authHasher)
	apiV1 := app.Group("/api/v1")
	apiV1.Use(bearerMiddleware(sessionRepository))
	apiV1.Post("/users", register.Register)
	apiV1.Post("/auth/login", login.Login)
	apiV1.Delete("/auth/logout", func(ctx *fiber.Ctx) error {
		return loginRequired(ctx, logout.Logout)
	})
	return &fibberApplication{app: app}
}

func bearerMiddleware(sessionRepository accountsrepos.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(requestcontext.Key, requestcontext.NewAnonymous())
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Next()
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		return validateSession(c, sessionRepository, token)
	}
}

func validateSession(c *fiber.Ctx, sessionRepository accountsrepos.SessionRepository, token string) error {
	validator := sessionvalidator.New(sessionRepository)
	session, err := validator.Validate(c.Context(), token)
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
	c.Locals(handler.SessionTokenKey, session.Token())
	return c.Next()
}

func loginRequired(ctx *fiber.Ctx, handlerFn func(*fiber.Ctx) error) error {
	requestCtx := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	if requestCtx.IsAuthenticated() {
		return handlerFn(ctx)
	}
	return ctx.SendStatus(http.StatusUnauthorized)
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
	defer func() { ctx.Status(code) }()
	ok := errors.As(err, &odinError)
	if ok {
		switch odinError.Tag() {
		case odinerrors.Domain:
			code = http.StatusBadRequest
		case odinerrors.NotFound:
			code = http.StatusNotFound
		case odinerrors.Unauthorized:
			code = http.StatusUnauthorized
		case odinerrors.AlreadyExists:
			code = http.StatusConflict
		default:
		}
		return ctx.JSON(map[string]string{"error": odinError.ExternalError()})
	}
	return nil
}
