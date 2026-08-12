package logouthandler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionterminator"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
)

type LogoutHandler interface {
	HandleResponse() error
}

type logoutHandler struct {
	sessionRepository repositories.SessionRepository
	handler           LogoutHandler
}

func New(sessionRepository repositories.SessionRepository, handler LogoutHandler) *logoutHandler {
	return &logoutHandler{
		sessionRepository: sessionRepository,
		handler:           handler,
	}
}

func (self *logoutHandler) Logout(ctx *fiber.Ctx) error {
	token := self.extractToken(ctx)
	terminator := sessionterminator.New(self.sessionRepository)
	if err := terminator.Terminate(ctx.Context(), token); err != nil {
		return err
	}
	return self.handler.HandleResponse()
}

func (self *logoutHandler) extractToken(ctx *fiber.Ctx) string {
	if cookie := ctx.Cookies(handler.SessionName); cookie != "" {
		return cookie
	}
	auth := ctx.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
