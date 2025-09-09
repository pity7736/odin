package htmxloginhandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
)

type HtmxLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) loginhandler.LoginHandler {
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
	self.ctx.Set("HX-Redirect", self.ctx.Query("next", "/"))
	return nil
}

func (self *HtmxLoginHandler) HandleBadRequest(err error) error {
	self.ctx.Render("login_error", LoginData{err.Error(), "/"}, "")
	return nil
}

func (self *HtmxLoginHandler) ContentType() string {
	return fiber.MIMETextHTMLCharsetUTF8
}

type LoginData struct {
	Error string
	Next  string
}
