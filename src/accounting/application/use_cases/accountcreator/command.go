package accountcreator

import (
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

type CreateAccountCommand struct {
	name           string
	initialBalance moneymodel.Money
	accountType    accounttypemodel.AccountType
}

func NewCreateAccountCommand(name string, initialBalance moneymodel.Money, accountType accounttypemodel.AccountType) CreateAccountCommand {
	return CreateAccountCommand{
		name:           name,
		initialBalance: initialBalance,
		accountType:    accountType,
	}
}

func (self CreateAccountCommand) Name() string {
	return self.name
}

func (self CreateAccountCommand) InitialBalance() moneymodel.Money {
	return self.initialBalance
}

func (self CreateAccountCommand) AccountType() accounttypemodel.AccountType {
	return self.accountType
}
