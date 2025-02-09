package loginhandler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/infrastructure/accountsrepositoryfactory"
)

type LoginHandler interface {
	HandleResponse(session *sessionmodel.Session) error
	HandleBadRequest(err error) error
	ContentType() string
}

type loginHandler struct {
	factory accountsrepositoryfactory.AccountsRepositoryFactory
	handler LoginHandler
}

func New(factory accountsrepositoryfactory.AccountsRepositoryFactory, handler LoginHandler) *loginHandler {
	return &loginHandler{factory: factory, handler: handler}
}

func (self *loginHandler) Login(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", self.handler.ContentType())
	var body LoginBody
	if err := self.validateRequestBody(ctx, &body); err != nil {
		ctx.Status(http.StatusBadRequest)
		return self.handler.HandleBadRequest(err)
	}
	return self.login(ctx, &body)
}

func (self *loginHandler) validateRequestBody(ctx *fiber.Ctx, body *LoginBody) error {
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

func (self *loginHandler) login(ctx *fiber.Ctx, body *LoginBody) error {
	sessionStarter := sessionstarter.New(body.Email, body.Password, self.factory)
	session, err := sessionStarter.Start(ctx.Context())
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return self.handler.HandleBadRequest(err)
	}
	ctx.Status(http.StatusCreated)
	return self.handler.HandleResponse(session)
}
