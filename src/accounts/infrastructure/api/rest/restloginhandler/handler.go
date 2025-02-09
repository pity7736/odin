package restloginhandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
)

type RestLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) loginhandler.LoginHandler {
	return &RestLoginHandler{ctx: ctx}
}

func (self *RestLoginHandler) HandleResponse(session *sessionmodel.Session) error {
	self.ctx.JSON(response{Token: session.Token(), Error: ""})
	return nil
}

func (self *RestLoginHandler) HandleBadRequest(err error) error {
	self.ctx.JSON(response{Token: "", Error: err.Error()})
	return nil
}

func (self *RestLoginHandler) ContentType() string {
	return fiber.MIMEApplicationJSON
}

type response struct {
	Token string `json:"token"`
	Error string `json:"error"`
}
