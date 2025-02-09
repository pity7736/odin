package loginhandler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/infrastructure/accountsrepositoryfactory"
)

type Handler struct {
	factory accountsrepositoryfactory.AccountsRepositoryFactory
}

func New(factory accountsrepositoryfactory.AccountsRepositoryFactory) *Handler {
	return &Handler{factory: factory}
}

func (self *Handler) Login(ctx *fiber.Ctx) error {
	var body LoginBody
	if err := self.validateRequestBody(ctx, &body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(response{Token: "", Error: err.Error()})
		return nil
	}
	return self.login(ctx, &body)
}

func (self *Handler) validateRequestBody(ctx *fiber.Ctx, body *LoginBody) error {
	if err := ctx.BodyParser(body); err != nil {
		return errors.New("wrong body")
	}
	if body.Email == "" {
		return errors.New("email is required")
	}
	if body.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (self *Handler) login(ctx *fiber.Ctx, body *LoginBody) error {
	sessionStarter := sessionstarter.New(body.Email, body.Password, self.factory)
	session, err := sessionStarter.Start(ctx.Context())
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(response{Token: "", Error: err.Error()})
		return nil
	}
	ctx.JSON(response{Token: session.Token(), Error: ""})
	ctx.Status(http.StatusCreated)
	return nil
}

type response struct {
	Token string `json:"token"`
	Error string `json:"error"`
}
