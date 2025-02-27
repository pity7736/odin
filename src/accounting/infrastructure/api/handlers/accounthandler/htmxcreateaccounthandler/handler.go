package htmxcreateaccounthandler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type HTMXCreateAccountHandler struct {
	repository repositories.AccountRepository
}

func New(repository repositories.AccountRepository) HTMXCreateAccountHandler {
	return HTMXCreateAccountHandler{repository: repository}
}

func (self HTMXCreateAccountHandler) Handle(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
	var body createAccountBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return errors.New("body is not valid")
	}
	var initialBalance moneymodel.Money
	if body.Name == "" {
		err = errors.New("name cannot be empty")
	} else {
		initialBalance, err = moneymodel.New(body.RawInitialBalance)
	}
	if err != nil {
		ctx.Render("create_account_error", RequestError{err.Error()}, "")
		return err
	}
	accountCreator := accountcreator.New(
		body.Name,
		ctx.Locals("userID").(string),
		initialBalance,
		self.repository,
	)
	account, _ := accountCreator.Create(context.TODO())
	return ctx.Render("account_created", account, "")
}

type RequestError struct {
	Error string
}

type createAccountBody struct {
	Name              string `form:"name"`
	RawInitialBalance string `form:"initial_balance"`
}
