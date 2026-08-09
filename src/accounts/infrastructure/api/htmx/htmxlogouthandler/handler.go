package htmxlogouthandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/infrastructure/api/logouthandler"
)

type htmxLogoutHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) logouthandler.LogoutHandler {
	return &htmxLogoutHandler{ctx: ctx}
}

func (self *htmxLogoutHandler) HandleResponse() error {
	self.ctx.Cookie(&fiber.Cookie{
		Name:     "__Secure-odin-session",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
	})
	self.ctx.Set("HX-Redirect", "/auth/login")
	return nil
}
