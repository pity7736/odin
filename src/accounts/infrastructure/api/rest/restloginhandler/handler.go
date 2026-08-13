package restloginhandler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type RestLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) loginhandler.LoginHandler {
	return &RestLoginHandler{ctx: ctx}
}

func (self *RestLoginHandler) HandleResponse(session *sessionmodel.Session) error {
	return self.ctx.JSON(response{Token: session.Token(), Error: ""})
}

func (self *RestLoginHandler) HandleBadRequest(err error) error {
	var odinError *odinerrors.Error
	errors.As(err, &odinError)
	return self.ctx.JSON(response{Token: "", Error: odinError.ExternalError()})
}

func (self *RestLoginHandler) ContentType() string {
	return fiber.MIMEApplicationJSON
}

type response struct {
	Token string `json:"token"`
	Error string `json:"error"`
}
