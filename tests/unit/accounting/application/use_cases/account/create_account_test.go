package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/application/use_cases/accountcreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
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
		repository.EXPECT().Add(mock.IsType(context.TODO()), mock.Anything).Return(nil)
		accountCreator := accountcreator.New(accountName, user.ID(), initialBalance, repository)

		account, err := accountCreator.Create(context.TODO())

		assert.Equal(t, user.ID(), account.UserID())
		assert.Equal(t, accountName, account.Name())
		assert.Equal(t, initialBalance, account.InitialBalance())
		assert.Equal(t, initialBalance, account.Balance())
		assert.True(t, testutils.IsUUIDv7(account.ID()))
		assert.Nil(t, err)
	})

	t.Run("should return error when initial balance is negative", func(t *testing.T) {
		user := userbuilder.New().Build()
		accountName := "saving account"
		initialBalance, _ := moneymodel.New("-1000000")
		repository := mocks.NewMockAccountRepository(t)
		accountCreator := accountcreator.New(accountName, user.ID(), initialBalance, repository)

		account, err := accountCreator.Create(context.TODO())

		assert.Nil(t, account)
		assert.Equal(t, errors.New("initial balance must be positive"), err)
		repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		user := userbuilder.New().Build()
		accountName := "saving account"
		initialBalance, _ := moneymodel.New("1000000")
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.IsType(context.TODO()), mock.Anything).Return(errors.New("some error"))
		accountCreator := accountcreator.New(accountName, user.ID(), initialBalance, repository)

		account, err := accountCreator.Create(context.TODO())

		assert.Nil(t, account)
		assert.Equal(t, errors.New("some error"), err)
	})
}
