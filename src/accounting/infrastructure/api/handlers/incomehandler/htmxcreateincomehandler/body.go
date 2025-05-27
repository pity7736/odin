package htmxcreateincomehandler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/incomecreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type createIncomeBody struct {
	Amount     string `json:"amount"`
	Date       string `json:"date"`
	CategoryID string `json:"category_id"`
	accountID  string
}

func newFromCtx(ctx *fiber.Ctx) (createIncomeBody, error) {
	accountID := ctx.Params("accountID")
	if accountID == "" {
		return createIncomeBody{}, odinerrors.NewErrorBuilder("account id is not sent").
			WithExternalMessage("el id de la cuenta es requerido").
			Build()
	}
	var body createIncomeBody
	if err := ctx.BodyParser(&body); err != nil {
		return createIncomeBody{}, err
	}
	body.accountID = accountID
	return body, nil
}

func (self createIncomeBody) toCommand() (incomecreator.CreateIncomeCommand, error) {
	amount, _ := moneymodel.New(self.Amount)
	date, _ := time.Parse(time.DateOnly, self.Date)
	return incomecreator.NewCommand(amount, date, self.CategoryID, self.accountID), nil
}
