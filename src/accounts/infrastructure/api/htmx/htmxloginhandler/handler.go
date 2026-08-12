package htmxloginhandler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
)

type HtmxLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) loginhandler.LoginHandler {
	return &HtmxLoginHandler{ctx: ctx}
}

func (self *HtmxLoginHandler) HandleResponse(session *sessionmodel.Session) error {
	cookie := fiber.Cookie{
		Name:     handler.SessionName,
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
	var odinError *odinerrors.Error
	errors.As(err, &odinError)
	return self.ctx.Render("login_error", LoginData{Error: odinError.ExternalError(), Next: "/"}, "")
}

func (self *HtmxLoginHandler) ContentType() string {
	return fiber.MIMETextHTMLCharsetUTF8
}

type LoginData struct {
	Error string
	Next  string
}
