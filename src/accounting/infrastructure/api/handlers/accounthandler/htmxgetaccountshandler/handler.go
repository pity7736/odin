package htmxgetaccountshandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type HTMXGetAccountsHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) HTMXGetAccountsHandler {
	return HTMXGetAccountsHandler{repository: repository}
}

func (self HTMXGetAccountsHandler) Handle(ctx *fiber.Ctx) error {
	requestContext, _ := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	accounts, err := self.repository.GetAll(context.WithValue(ctx.Context(), requestcontext.Key, requestContext))
	if err != nil {
		return err
	}
	return ctx.Render("accounts", data{Accounts: accounts})
}

type data struct {
	Accounts []*accountmodel.Account
}
