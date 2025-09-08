package pgrepositories

import (
	"context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
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

func (self *PGAccountRepository) GetByID(ctx context.Context, id string) (*accountmodel.Account, error) {
	requestContext := ctx.Value(requestcontext.Key).(*requestcontext.RequestContext)
	for _, account := range self.accounts {
		if account.ID() == id && account.UserID() == requestContext.UserID() {
			return account, nil
		}
	}
	return nil, odinerrors.NewErrorBuilder("account not found").WithExternalMessage("account not found").Build()
}

func (self *PGAccountRepository) Save(ctx context.Context, account *accountmodel.Account) error {
	//TODO implement me
	panic("implement me")
}
