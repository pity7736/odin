package createaccounthandler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type CreateAccountHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) *CreateAccountHandler {
	return &CreateAccountHandler{repository: repository}
}

func (self *CreateAccountHandler) Handle(ctx *fiber.Ctx) (*accountmodel.Account, error) {
	command, err := self.createCommand(ctx)
	if err != nil {
		return nil, err
	}
	accountCreator := accountcreator.New(*command, self.repository)
	return accountCreator.Create(context.TODO())
}

func (self *CreateAccountHandler) createCommand(ctx *fiber.Ctx) (*accountcreator.CreateAccountCommand, error) {
	var body createAccountBody
	if err := ctx.BodyParser(&body); err != nil {
		return nil, err
	}
	return body.toCommand(ctx.Locals("userID").(string))
}

type createAccountBody struct {
	Name              string `json:"name" form:"name"`
	RawInitialBalance string `json:"initial_balance" form:"initial_balance"`
}

func (self createAccountBody) toCommand(userID string) (*accountcreator.CreateAccountCommand, error) {
	initialBalance, err := moneymodel.New(self.RawInitialBalance)
	if err != nil {
		return nil, err
	}
	return accountcreator.NewCreateAccountCommand(self.Name, initialBalance, userID), err
}
