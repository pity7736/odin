package account_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/mocks"
)

func TestAccountCreator(t *testing.T) {

	t.Run("should create an account", func(t *testing.T) {
		user := userbuilder.New().Build()
		accountName := "saving account"
		initialBalance, _ := moneymodel.New("1000000")
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		command := accountcreator.NewCreateAccountCommand(accountName, initialBalance)
		accountCreator := accountcreator.New(*command, repository)
		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)

		account, err := accountCreator.Create(ctx)

		assert.Equal(t, user.ID(), account.UserID())
		assert.Equal(t, accountName, account.Name())
		assert.Equal(t, initialBalance, account.InitialBalance())
		assert.Equal(t, initialBalance, account.Balance())
		assert.True(t, testutils.IsUUIDv7(account.ID()))
		assert.True(t, testutils.IsTimeClose(time.Now(), account.CreatedAt()))
		assert.Nil(t, err)
		repository.AssertCalled(t, "Add", ctx, account)
	})

	t.Run("should return error when initial balance is negative", func(t *testing.T) {
		user := userbuilder.New().Build()
		accountName := "saving account"
		initialBalance, _ := moneymodel.New("-1000000")
		repository := mocks.NewMockAccountRepository(t)
		command := accountcreator.NewCreateAccountCommand(accountName, initialBalance)
		accountCreator := accountcreator.New(*command, repository)
		requestContext, _ := requestcontext.New(user.ID())

		account, err := accountCreator.Create(context.WithValue(context.TODO(), requestcontext.Key, requestContext))

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Nil(t, account)
		assert.Equal(t, "validation error: initial balance must be positive", odinError.ExternalError())
		repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		user := userbuilder.New().Build()
		accountName := "saving account"
		initialBalance, _ := moneymodel.New("1000000")
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(errors.New("some error"))
		command := accountcreator.NewCreateAccountCommand(accountName, initialBalance)
		accountCreator := accountcreator.New(*command, repository)

		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		account, err := accountCreator.Create(ctx)

		assert.Nil(t, account)
		assert.Equal(t, errors.New("some error"), err)
		repository.AssertCalled(t, "Add", ctx, mock.Anything)
	})
}
