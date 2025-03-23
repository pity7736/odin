package htmxcreateaccounthandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/createaccounthandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
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
			return odinerrors.NewErrorBuilder("error rendering create account error block").
				WithWrapped(renderError).
				WithTag(odinerrors.RENDER).
				Build()
		}
		return odinerrors.NewErrorBuilder("error creating account").
			WithWrapped(err).
			WithExternalMessage("error creating account").
			Build()
	}
	renderError := ctx.Render("account_created", account, "")
	if renderError != nil {
		return odinerrors.NewErrorBuilder("error rendering create account block").
			WithWrapped(renderError).
			WithTag(odinerrors.RENDER).
			Build()
	}
	return nil
}
