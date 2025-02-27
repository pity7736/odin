package accountcreator

import moneymodel "raiseexception.dev/odin/src/accounting/domain/money"

type CreateAccountCommand struct {
	name           string
	initialBalance moneymodel.Money
	userID         string
}

func NewCreateAccountCommand(name string, initialBalance moneymodel.Money, userID string) *CreateAccountCommand {
	return &CreateAccountCommand{
		name:           name,
		initialBalance: initialBalance,
		userID:         userID,
	}
}

func (self *CreateAccountCommand) Name() string {
	return self.name
}

func (self *CreateAccountCommand) InitialBalance() moneymodel.Money {
	return self.initialBalance
}

func (self *CreateAccountCommand) UserID() string {
	return self.userID
}
