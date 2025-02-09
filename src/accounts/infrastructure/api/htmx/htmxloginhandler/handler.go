package htmxloginhandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
)

type HtmxLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) *HtmxLoginHandler {
	return &HtmxLoginHandler{ctx: ctx}
}

func (self *HtmxLoginHandler) HandleResponse(session *sessionmodel.Session) error {
	cookie := fiber.Cookie{
		Name:     "__Secure-odin-session",
		Value:    session.Token(),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
	}
	self.ctx.Cookie(&cookie)
	self.ctx.Set("HX-Redirect", "/")
	return nil
}

func (self *HtmxLoginHandler) HandleBadRequest(err error) error {
	self.ctx.Render("login_error", RequestError{err.Error()}, "")
	return nil
}

func (self *HtmxLoginHandler) ContentType() string {
	return fiber.MIMETextHTMLCharsetUTF8
}

type RequestError struct {
	Error string
}
