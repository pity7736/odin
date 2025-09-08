package htmxgetaccounthandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type HTMXGetAccountHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) HTMXGetAccountHandler {
	return HTMXGetAccountHandler{repository: repository}
}

func (self HTMXGetAccountHandler) Handle(ctx *fiber.Ctx) error {
	requestContext, _ := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	account, err := self.repository.GetByID(
		context.WithValue(ctx.Context(), requestcontext.Key, requestContext),
		ctx.Params("accountID"),
	)
	if err != nil {
		return err
	}
	return ctx.Render("account", account, "")
}
