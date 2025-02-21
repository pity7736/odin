package accountmodel

import moneymodel "raiseexception.dev/odin/src/accounting/domain/money"

type Account struct {
	name           string
	initialBalance moneymodel.Money
	userID         string
	id             string
	balance        moneymodel.Money
}

func New(id, name, userID string, initialBalance, balance moneymodel.Money) *Account {
	return &Account{
		id:             id,
		name:           name,
		initialBalance: initialBalance,
		userID:         userID,
		balance:        balance,
	}
}

func (self *Account) ID() string {
	return self.id
}

func (self *Account) Name() string {
	return self.name
}

func (self *Account) InitialBalance() moneymodel.Money {
	return self.initialBalance
}

func (self *Account) UserID() string {
	return self.userID
}

func (self *Account) Balance() moneymodel.Money {
	return self.balance
}
