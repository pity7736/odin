package accountcreator

import (
	"context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type AccountCreator struct {
	command    CreateAccountCommand
	repository repositories.AccountRepository
}

func New(command CreateAccountCommand, repository repositories.AccountRepository) *AccountCreator {
	return &AccountCreator{command: command, repository: repository}
}

func (self *AccountCreator) Create(ctx context.Context) (*accountmodel.Account, error) {
	requestContext := ctx.Value(requestcontext.Key).(*requestcontext.RequestContext)
	currency := self.command.InitialBalance().Currency()
	exists, err := self.repository.ExistsByNameAndCurrency(ctx, self.command.Name(), currency)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, odinerrors.NewErrorBuilder("account with name and currency already exists for user").
			WithExternalMessage("Ya tienes una cuenta con ese nombre en esa moneda").
			WithTag(odinerrors.Domain).
			Build()
	}
	account, err := accountmodel.New(self.command.Name(), requestContext.UserID(), self.command.InitialBalance(), self.command.AccountType())
	if err != nil {
		return nil, err
	}
	err = self.repository.Add(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
