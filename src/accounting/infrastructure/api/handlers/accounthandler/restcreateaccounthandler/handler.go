package restcreateaccounthandler

import (
	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/createaccounthandler"
)

type RestCreateAccountHandler struct {
	handler *createaccounthandler.CreateAccountHandler
}

func New(repository repositories.AccountRepository) RestCreateAccountHandler {
	return RestCreateAccountHandler{handler: createaccounthandler.New(repository)}
}

func (self RestCreateAccountHandler) Handle(ctx *fiber.Ctx) error {
	ctx.Set("content-type", fiber.MIMEApplicationJSON)
	account, err := self.handler.Handle(ctx)
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

type createAccountResponse struct {
	Name           string `json:"name"`
	InitialBalance string `json:"initial_balance"`
	Balance        string `json:"balance"`
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
}
