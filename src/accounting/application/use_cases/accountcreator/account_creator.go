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
	account, err := accountmodel.New(self.command.Name(), requestContext.UserID(), self.command.InitialBalance())
	if err != nil {
		return nil, odinerrors.NewErrorBuilder("error creating a new account").
			WithWrapped(err).
			WithExternalMessage("validation error").
			Build()
	}
	err = self.repository.Add(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
