package htmxcreateaccounthandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/createaccounthandler"
)

type HTMXCreateAccountHandler struct {
	handler *createaccounthandler.CreateAccountHandler
}

func New(repository repositories.AccountRepository) HTMXCreateAccountHandler {
	return HTMXCreateAccountHandler{handler: createaccounthandler.New(repository)}
}

func (self HTMXCreateAccountHandler) Handle(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
	account, err := self.handler.Handle(ctx)
	if err != nil {
		renderError := ctx.Render("create_account_error", err, "")
		if renderError != nil {
			return renderError
		}
		return err
	}
	return ctx.Render("account_created", account, "")
}
