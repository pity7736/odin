package pgrepositories

import (
	"context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
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

func (self *PGAccountRepository) GetAll(ctx context.Context) ([]*accountmodel.Account, error) {
	requestContext := ctx.Value(requestcontext.Key).(*requestcontext.RequestContext)
	result := make([]*accountmodel.Account, 0, len(self.accounts))
	for _, account := range self.accounts {
		if account.UserID() == requestContext.UserID() {
			result = append(result, account)
		}
	}
	return result, nil
}
