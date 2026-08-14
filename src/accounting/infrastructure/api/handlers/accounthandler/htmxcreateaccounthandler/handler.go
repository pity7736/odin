package htmxcreateaccounthandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/accountviewmodel"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/createaccounthandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
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
		renderErr := ctx.Render("create_account_error", fiber.Map{"ExternalError": handler.ExternalOrFallback(err, "No se pudo crear la cuenta")}, "")
		if renderErr != nil {
			return odinerrors.NewErrorBuilder("error rendering create account error block").
				WithWrapped(renderErr).
				WithTag(odinerrors.Render).
				Build()
		}
		return err
	}
	renderErr := ctx.Render("account_created_oob", accountviewmodel.New(account), "")
	if renderErr != nil {
		return odinerrors.NewErrorBuilder("error rendering create account block").
			WithWrapped(renderErr).
			WithTag(odinerrors.Render).
			Build()
	}
	return nil
}
