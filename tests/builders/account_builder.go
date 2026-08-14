package builders

import (
	"context"
	"time"

	"github.com/google/uuid"

	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type AccountBuilder struct {
	name           string
	initialBalance moneymodel.Money
	balance        moneymodel.Money
	userID         string
	id             string
	accountType    accounttypemodel.AccountType
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
		accountType:    accounttypemodel.Savings(),
	}
}

func (self *AccountBuilder) WithUserID(userID string) *AccountBuilder {
	self.userID = userID
	return self
}

func (self *AccountBuilder) WithName(name string) *AccountBuilder {
	self.name = name
	return self
}

func (self *AccountBuilder) WithInitialBalance(value string) *AccountBuilder {
	currency := self.initialBalance.Currency()
	balance, _ := moneymodel.New(value, currency)
	self.initialBalance = balance
	self.balance = balance
	return self
}

func (self *AccountBuilder) WithType(accountType accounttypemodel.AccountType) *AccountBuilder {
	self.accountType = accountType
	return self
}

func (self *AccountBuilder) WithCurrency(currency moneymodel.Currency) *AccountBuilder {
	self.initialBalance = moneymodel.MustNew(self.initialBalance.String(), currency)
	self.balance = moneymodel.MustNew(self.balance.String(), currency)
	return self
}

func (self *AccountBuilder) Build() *accountmodel.Account {
	account, _ := accountmodel.NewFromRepository(self.id, self.name, self.userID, self.initialBalance, self.balance, self.accountType, time.Now())
	return account
}

func (self *AccountBuilder) Create(repository repositories.AccountRepository) *accountmodel.Account {
	requestContext, _ := requestcontext.New(self.userID)
	ctx := context.WithValue(context.Background(), requestcontext.Key, requestContext)
	accountCreator := accountcreator.New(
		accountcreator.NewCreateAccountCommand(self.name, self.initialBalance, self.accountType),
		repository,
	)
	account, _ := accountCreator.Create(ctx)
	return account
}
