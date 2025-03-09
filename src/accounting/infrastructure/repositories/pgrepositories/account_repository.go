package pgrepositories

import (
	"context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
)

type PGAccountRepository struct {
	accounts map[string]*accountmodel.Account
}

func NewAccountRepository() *PGAccountRepository {
	return &PGAccountRepository{accounts: make(map[string]*accountmodel.Account)}
}

func (self *PGAccountRepository) Add(ctx context.Context, account *accountmodel.Account) error {
	self.accounts[account.ID()] = account
	return nil
}

func (self *PGAccountRepository) GetAll(ctx context.Context, userID string) []*accountmodel.Account {
	result := make([]*accountmodel.Account, 0, len(self.accounts))
	for _, account := range self.accounts {
		if account.UserID() == userID {
			result = append(result, account)
		}
	}
	return result
}
