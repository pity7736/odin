package loginhandler

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/infrastructure/accountsrepositoryfactory"
)

type Handler struct {
	factory accountsrepositoryfactory.AccountsRepositoryFactory
}

func New(factory accountsrepositoryfactory.AccountsRepositoryFactory) *Handler {
	return &Handler{factory: factory}
}

func (self *Handler) Login(ctx *fiber.Ctx) error {
	var body loginBody
	if err := ctx.BodyParser(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(fiber.Map{"error": "wrong body"})
		return nil
	}
	if body.Email == "" {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(fiber.Map{"error": "email is required"})
		return nil
	}
	if body.Password == "" {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(fiber.Map{"error": "password is required"})
		return nil
	}
	sessionStarter := sessionstarter.New(body.Email, body.Password, self.factory)
	_, err := sessionStarter.Start(ctx.Context())
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.JSON(fiber.Map{"error": err.Error()})
		return nil
	}
	return nil
}
