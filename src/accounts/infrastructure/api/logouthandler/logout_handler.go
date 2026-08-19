package logouthandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionterminator"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
)

type logoutHandler struct {
	sessionRepository repositories.SessionRepository
	ctx               *fiber.Ctx
}

func New(sessionRepository repositories.SessionRepository, ctx *fiber.Ctx) *logoutHandler {
	return &logoutHandler{
		sessionRepository: sessionRepository,
		ctx:               ctx,
	}
}

func (self *logoutHandler) Logout() error {
	token, _ := self.ctx.Locals(handler.SessionTokenKey).(string)
	terminator := sessionterminator.New(self.sessionRepository)
	if err := terminator.Terminate(self.ctx.Context(), token); err != nil {
		return err
	}
	return self.ctx.JSON(map[string]string{"message": "session closed"})
}
