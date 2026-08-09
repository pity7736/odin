package restlogouthandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/infrastructure/api/logouthandler"
)

type restLogoutHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) logouthandler.LogoutHandler {
	return &restLogoutHandler{ctx: ctx}
}

func (self *restLogoutHandler) HandleResponse() error {
	return self.ctx.JSON(map[string]string{"message": "session closed"})
}
