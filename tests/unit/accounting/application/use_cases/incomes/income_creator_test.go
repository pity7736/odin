package incomes_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/application/use_cases/incomecreator"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestIncomeCreatorShould(t *testing.T) {

	t.Run("return error when account does not exist", func(t *testing.T) {
		accountingFactory := testrepositoryfactory.New(t)
		userRepositoryMock := accountingFactory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepositoryMock)
		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		accountID := "1234"
		accountError := odinerrors.NewErrorBuilder("account not found").
			WithTag(odinerrors.NOT_FOUND).
			WithExternalMessage(fmt.Sprintf("account with id %s does not exist", accountID)).
			Build()
		accountRepositoryMock := accountingFactory.GetAccountRepositoryMock()
		accountRepositoryMock.EXPECT().GetByID(ctx, accountID).Return(nil, accountError)

		categoryRepositoryMock := accountingFactory.GetCategoryRepositoryMock()
		categoryRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithUser(user).Create(categoryRepositoryMock)
		//categoryRepositoryMock.EXPECT().GetByID(ctx, category.ID()).Return(category, nil)
		amount, _ := moneymodel.New("100000")
		incomeCreator := incomecreator.New(
			accountingFactory,
			amount,
			time.Now(),
			category.ID(),
			accountID,
		)
		incomeRepositoryMock := accountingFactory.GetIncomeRepositoryMock()

		income, err := incomeCreator.Create(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)

		assert.Nil(t, income)
		assert.True(t, ok)
		assert.Equal(t, odinerrors.NOT_FOUND, odinError.Tag())
		assert.Equal(t, "account with id 1234 does not exist", odinError.ExternalError())
		incomeRepositoryMock.AssertNotCalled(t, "Add")
	})

	t.Run("return error when category does not exist", func(t *testing.T) {
		accountingFactory := testrepositoryfactory.New(t)
		userRepositoryMock := accountingFactory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepositoryMock)

		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		categoryID := "1234"
		categoryError := odinerrors.NewErrorBuilder("account not found").
			WithTag(odinerrors.NOT_FOUND).
			WithExternalMessage(fmt.Sprintf("category with id %s does not exist", categoryID)).
			Build()
		categoryRepositoryMock := accountingFactory.GetCategoryRepositoryMock()
		categoryRepositoryMock.EXPECT().GetByID(ctx, categoryID).Return(nil, categoryError)

		accountRepositoryMock := accountingFactory.GetAccountRepositoryMock()
		accountRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepositoryMock)
		accountRepositoryMock.EXPECT().GetByID(ctx, account.ID()).Return(account, nil)

		amount, _ := moneymodel.New("100000")
		incomeCreator := incomecreator.New(
			accountingFactory,
			amount,
			time.Now(),
			categoryID,
			account.ID(),
		)
		incomeRepositoryMock := accountingFactory.GetIncomeRepositoryMock()

		income, err := incomeCreator.Create(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)

		assert.Nil(t, income)
		assert.True(t, ok)
		assert.Equal(t, odinerrors.NOT_FOUND, odinError.Tag())
		assert.Equal(t, "category with id 1234 does not exist", odinError.ExternalError())
		incomeRepositoryMock.AssertNotCalled(t, "Add")
	})

	t.Run("return error when category does not belong to user", func(t *testing.T) {
		accountingFactory := testrepositoryfactory.New(t)
		userRepositoryMock := accountingFactory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepositoryMock)

		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		categoryRepositoryMock := accountingFactory.GetCategoryRepositoryMock()
		categoryRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().Create(categoryRepositoryMock)
		categoryRepositoryMock.EXPECT().GetByID(ctx, category.ID()).Return(category, nil)

		accountRepositoryMock := accountingFactory.GetAccountRepositoryMock()
		accountRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepositoryMock)
		accountRepositoryMock.EXPECT().GetByID(ctx, account.ID()).Return(account, nil)

		amount, _ := moneymodel.New("100000")
		incomeCreator := incomecreator.New(
			accountingFactory,
			amount,
			time.Now(),
			category.ID(),
			account.ID(),
		)
		incomeRepositoryMock := accountingFactory.GetIncomeRepositoryMock()

		income, err := incomeCreator.Create(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)

		assert.Nil(t, income)
		assert.True(t, ok)
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
		assert.Equal(t, "la categoría no pertenece al usuario logueado", odinError.ExternalError())
		incomeRepositoryMock.AssertNotCalled(t, "Add")
	})

	t.Run("return error when account does not belong to user", func(t *testing.T) {
		accountingFactory := testrepositoryfactory.New(t)
		userRepositoryMock := accountingFactory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepositoryMock)

		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		categoryRepositoryMock := accountingFactory.GetCategoryRepositoryMock()
		categoryRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithUser(user).Create(categoryRepositoryMock)
		//categoryRepositoryMock.EXPECT().GetByID(ctx, category.ID()).Return(category, nil)

		accountRepositoryMock := accountingFactory.GetAccountRepositoryMock()
		accountRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().Create(accountRepositoryMock)
		accountRepositoryMock.EXPECT().GetByID(ctx, account.ID()).Return(account, nil)

		amount, _ := moneymodel.New("100000")
		incomeCreator := incomecreator.New(
			accountingFactory,
			amount,
			time.Now(),
			category.ID(),
			account.ID(),
		)
		incomeRepositoryMock := accountingFactory.GetIncomeRepositoryMock()

		income, err := incomeCreator.Create(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)

		assert.Nil(t, income)
		assert.True(t, ok)
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
		assert.Equal(t, "la cuenta no pertenece al usuario logueado", odinError.ExternalError())
		incomeRepositoryMock.AssertNotCalled(t, "Add")
	})

	t.Run("return income when data is valid", func(t *testing.T) {
		accountingFactory := testrepositoryfactory.New(t)
		userRepositoryMock := accountingFactory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepositoryMock)

		requestContext, _ := requestcontext.New(user.ID())
		ctx := context.WithValue(context.TODO(), requestcontext.Key, requestContext)
		categoryRepositoryMock := accountingFactory.GetCategoryRepositoryMock()
		categoryRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithUser(user).WithIncomeType().Create(categoryRepositoryMock)
		categoryRepositoryMock.EXPECT().GetByID(ctx, category.ID()).Return(category, nil)

		accountRepositoryMock := accountingFactory.GetAccountRepositoryMock()
		accountRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepositoryMock)
		accountRepositoryMock.EXPECT().GetByID(ctx, account.ID()).Return(account, nil)
		accountRepositoryMock.EXPECT().Save(ctx, account).Return(nil)

		incomeRepositoryMock := accountingFactory.GetIncomeRepositoryMock()
		incomeRepositoryMock.EXPECT().Add(ctx, mock.Anything).Return(nil)

		amount, _ := moneymodel.New("100000")
		incomeCreator := incomecreator.New(
			accountingFactory,
			amount,
			time.Now(),
			category.ID(),
			account.ID(),
		)

		income, err := incomeCreator.Create(ctx)

		assert.Nil(t, err)
		assert.Equal(t, amount, income.Amount())
		incomeRepositoryMock.AssertCalled(t, "Add", ctx, income)
	})
}

// TODO: add case when save account balance or add income fails. aka transactions
