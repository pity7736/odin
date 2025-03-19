package accountcreator

import (
	"context"
	"time"

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

func New(command CreateAccountCommand, repository repositories.AccountRepository) *AccountCreator {
	return &AccountCreator{
		name:           command.Name(),
		initialBalance: command.InitialBalance(),
		userID:         command.UserID(),
		repository:     repository,
	}
}

func (self *AccountCreator) Create(ctx context.Context) (*accountmodel.Account, error) {
	id, _ := uuid.NewV7()
	account, err := accountmodel.New(
		id.String(),
		self.name,
		self.userID,
		self.initialBalance,
		self.initialBalance,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	err = self.repository.Add(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
