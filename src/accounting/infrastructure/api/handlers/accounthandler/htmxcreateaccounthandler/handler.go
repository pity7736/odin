package htmxcreateaccounthandler

import (
	"errors"

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
		renderErr := ctx.Render("create_account_error", fiber.Map{"ExternalError": externalOrFallback(err)}, "")
		if renderErr != nil {
			return odinerrors.NewErrorBuilder("error rendering create account error block").
				WithWrapped(renderErr).
				WithTag(odinerrors.Render).
				Build()
		}
		return err
	}
	renderErr := ctx.Render("account_created", account, "")
	if renderErr != nil {
		return odinerrors.NewErrorBuilder("error rendering create account block").
			WithWrapped(renderErr).
			WithTag(odinerrors.Render).
			Build()
	}
	return nil
}

func externalOrFallback(err error) string {
	var odinError *odinerrors.Error
	if errors.As(err, &odinError) && odinError.ExternalError() != "" {
		return odinError.ExternalError()
	}
	return "No se pudo crear la cuenta"
}
