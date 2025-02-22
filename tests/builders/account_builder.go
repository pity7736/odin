package builders

import (
	"github.com/google/uuid"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

type AccountBuilder struct {
	name           string
	initialBalance moneymodel.Money
	balance        moneymodel.Money
	userID         string
	id             string
}

func NewAccountBuilder() *AccountBuilder {
	balance, _ := moneymodel.New("1000000")
	id, _ := uuid.NewV7()
	userID, _ := uuid.NewV7()
	return &AccountBuilder{
		name:           "test",
		initialBalance: balance,
		balance:        balance,
		userID:         userID.String(),
		id:             id.String(),
	}
}

func (self *AccountBuilder) WithUserID(userID string) *AccountBuilder {
	self.userID = userID
	return self
}

func (self *AccountBuilder) Build() *accountmodel.Account {
	return accountmodel.New(
		self.id,
		self.name,
		self.userID,
		self.initialBalance,
		self.balance,
	)
}
