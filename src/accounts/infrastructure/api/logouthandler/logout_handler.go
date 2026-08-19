package logouthandler

import (
	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionterminator"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
)

type logoutHandler struct {
	sessionRepository repositories.SessionRepository
}

func New(sessionRepository repositories.SessionRepository) logoutHandler {
	return logoutHandler{
		sessionRepository: sessionRepository,
	}
}

func (self logoutHandler) Logout(ctx *fiber.Ctx) error {
	token, _ := ctx.Locals(handler.SessionTokenKey).(string)
	terminator := sessionterminator.New(self.sessionRepository)
	if err := terminator.Terminate(ctx.Context(), token); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"message": "session closed"})
}
