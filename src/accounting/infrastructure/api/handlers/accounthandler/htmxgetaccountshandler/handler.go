package htmxgetaccountshandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type HTMXGetAccountsHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) HTMXGetAccountsHandler {
	return HTMXGetAccountsHandler{repository: repository}
}

func (self HTMXGetAccountsHandler) Handle(ctx *fiber.Ctx) error {
	accounts, err := self.repository.GetAll(context.TODO(), ctx.Locals("userID").(string))
	if err != nil {
		return err
	}
	return ctx.Render("accounts", data{Accounts: accounts})
}

type data struct {
	Accounts []*accountmodel.Account
}
