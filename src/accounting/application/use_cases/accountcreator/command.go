package accountcreator

import moneymodel "raiseexception.dev/odin/src/accounting/domain/money"

type CreateAccountCommand struct {
	name           string
	initialBalance moneymodel.Money
}

func NewCreateAccountCommand(name string, initialBalance moneymodel.Money) *CreateAccountCommand {
	return &CreateAccountCommand{
		name:           name,
		initialBalance: initialBalance,
	}
}

func (self *CreateAccountCommand) Name() string {
	return self.name
}

func (self *CreateAccountCommand) InitialBalance() moneymodel.Money {
	return self.initialBalance
}
