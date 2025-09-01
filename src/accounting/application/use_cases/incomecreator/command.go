package incomecreator

import (
	"time"

	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

type CreateIncomeCommand struct {
	amount     moneymodel.Money
	date       time.Time
	categoryID string
	accountID  string
}

func NewCommand(amount moneymodel.Money, date time.Time, categoryID, accountID string) CreateIncomeCommand {
	return CreateIncomeCommand{
		amount:     amount,
		date:       date,
		categoryID: categoryID,
		accountID:  accountID,
	}
}

func (self CreateIncomeCommand) Amount() moneymodel.Money {
	return self.amount
}

func (self CreateIncomeCommand) Date() time.Time {
	return self.date
}

func (self CreateIncomeCommand) CategoryID() string {
	return self.categoryID
}

func (self CreateIncomeCommand) AccountID() string {
	return self.accountID
}
