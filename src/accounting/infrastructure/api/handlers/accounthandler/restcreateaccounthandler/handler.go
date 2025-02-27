package restcreateaccounthandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type RestCreateAccountHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) RestCreateAccountHandler {
	return RestCreateAccountHandler{repository: repository}
}

func (self RestCreateAccountHandler) Handle(ctx *fiber.Ctx) error {
	ctx.Set("content-type", fiber.MIMEApplicationJSON)
	var body createAccountBody
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	initialBalance, err := moneymodel.New(body.RawInitialBalance)
	if err != nil {
		return err
	}
	accountCreator := accountcreator.New(
		body.Name,
		ctx.Locals("userID").(string),
		initialBalance,
		self.repository,
	)
	account, err := accountCreator.Create(context.TODO())
	if err != nil {
		return err
	}
	return ctx.JSON(createAccountResponse{
		Name:           account.Name(),
		ID:             account.ID(),
		InitialBalance: account.InitialBalance().String(),
		Balance:        account.Balance().String(),
		UserID:         account.UserID(),
	})
}

type createAccountBody struct {
	Name              string `json:"name"`
	RawInitialBalance string `json:"initial_balance"`
}

type createAccountResponse struct {
	Name           string `json:"name"`
	InitialBalance string `json:"initial_balance"`
	Balance        string `json:"balance"`
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
}
