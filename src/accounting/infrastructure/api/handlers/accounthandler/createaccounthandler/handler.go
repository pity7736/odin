package createaccounthandler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type CreateAccountHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) *CreateAccountHandler {
	return &CreateAccountHandler{repository: repository}
}

func (self *CreateAccountHandler) Handle(ctx *fiber.Ctx) (*accountmodel.Account, error) {
	requestContext := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	command, err := self.createCommand(ctx)
	if err != nil {
		return nil, err
	}
	accountCreator := accountcreator.New(command, self.repository)
	account, err := accountCreator.Create(context.WithValue(ctx.Context(), requestcontext.Key, requestContext))
	ctx.Status(fiber.StatusCreated)
	return account, err
}

func (self *CreateAccountHandler) createCommand(ctx *fiber.Ctx) (accountcreator.CreateAccountCommand, error) {
	var body createAccountBody
	if err := ctx.BodyParser(&body); err != nil {
		return accountcreator.CreateAccountCommand{}, odinerrors.NewErrorBuilder("wrong body").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			WithWrapped(err).
			Build()
	}
	return body.toCommand()
}

type createAccountBody struct {
	Name              string `json:"name" form:"name"`
	RawInitialBalance string `json:"initial_balance" form:"initial_balance"`
	RawType           string `json:"type" form:"type"`
	RawCurrency       string `json:"currency" form:"currency"`
}

func (self createAccountBody) toCommand() (accountcreator.CreateAccountCommand, error) {
	if self.RawInitialBalance == "" {
		return accountcreator.CreateAccountCommand{}, odinerrors.NewErrorBuilder("initial balance is required").
			WithExternalMessage("El saldo inicial es obligatorio").
			WithTag(odinerrors.Domain).
			Build()
	}
	currency, err := moneymodel.CurrencyFromString(strings.Clone(self.RawCurrency))
	if err != nil {
		return accountcreator.CreateAccountCommand{}, err
	}
	accountType, err := accounttypemodel.NewFromString(strings.Clone(self.RawType))
	if err != nil {
		return accountcreator.CreateAccountCommand{}, err
	}
	initialBalance, err := moneymodel.New(strings.Clone(self.RawInitialBalance), currency)
	if err != nil {
		return accountcreator.CreateAccountCommand{}, err
	}
	return accountcreator.NewCreateAccountCommand(strings.Clone(self.Name), initialBalance, accountType), nil
}
