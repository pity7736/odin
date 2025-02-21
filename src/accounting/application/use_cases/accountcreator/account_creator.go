package accountcreator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type AccountCreator struct {
	name           string
	initialBalance moneymodel.Money
	userID         string
	repository     repositories.AccountRepository
}

func New(name, userID string, initialBalance moneymodel.Money, repository repositories.AccountRepository) *AccountCreator {
	return &AccountCreator{
		name:           name,
		initialBalance: initialBalance,
		userID:         userID,
		repository:     repository,
	}
}

func (self *AccountCreator) Create(ctx context.Context) (*accountmodel.Account, error) {
	if self.initialBalance.IsNegative() {
		return nil, errors.New("initialBalance must be positive")
	}
	id, _ := uuid.NewV7()
	account := accountmodel.New(
		id.String(),
		self.name,
		self.userID,
		self.initialBalance,
		self.initialBalance,
	)
	err := self.repository.Add(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
