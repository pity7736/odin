package accountcreator

import (
	"context"

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

func New(command CreateAccountCommand, repository repositories.AccountRepository) *AccountCreator {
	return &AccountCreator{
		name:           command.Name(),
		initialBalance: command.InitialBalance(),
		userID:         command.UserID(),
		repository:     repository,
	}
}

func (self *AccountCreator) Create(ctx context.Context) (*accountmodel.Account, error) {
	account, err := accountmodel.New(self.name, self.userID, self.initialBalance)
	if err != nil {
		return nil, err
	}
	err = self.repository.Add(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
