package incomehandlers_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/incomehandler/htmxcreateincomehandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestHTMXCreateIncomeHandlerShould(t *testing.T) {

	t.Run("return error when account id is not sent", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetIncomeRepositoryMock()
		categoryRepository := factory.GetCategoryRepositoryMock()
		categoryRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().Create(categoryRepository)
		ctxBuilder := builders.NewFiberContextBuilder().
			WithMethod("POST").
			WithContentType(fiber.MIMEApplicationForm).
			WithBody(fmt.Appendf(nil,
				"amount=%s&date=%s&category_id=%s",
				"1000",
				"2025-04-03",
				category.ID(),
			))
		defer ctxBuilder.Release()
		ctx := ctxBuilder.Build()
		patches := gomonkey.ApplyMethodReturn(ctx, "Params", "")
		defer patches.Reset()
		createIncomeHandler := htmxcreateincomehandler.New(factory)

		err := createIncomeHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "el id de la cuenta es requerido", odinError.ExternalError())
		assert.True(t, strings.Contains(responseBody, odinError.ExternalError()))
		repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("return error when account does not exist", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetIncomeRepositoryMock()
		accountID := "12345"
		accountRepository := factory.GetAccountRepositoryMock()
		accountError := odinerrors.NewErrorBuilder("account not found").
			WithExternalMessage(fmt.Sprintf("no existe una cuenta con el id %s", accountID)).
			Build()
		accountRepository.EXPECT().GetByID(mock.Anything, accountID).Return(nil, accountError)
		categoryRepository := factory.GetCategoryRepositoryMock()
		categoryRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().Create(categoryRepository)
		ctxBuilder := builders.NewFiberContextBuilder().
			WithMethod("POST").
			WithContentType(fiber.MIMEApplicationForm).
			WithBody(fmt.Appendf(nil,
				"amount=%s&date=%s&category_id=%s",
				"1000",
				"2025-04-03",
				category.ID(),
			))
		defer ctxBuilder.Release()
		ctx := ctxBuilder.Build()

		patches := gomonkey.ApplyMethodReturn(ctx, "Params", accountID)
		defer patches.Reset()
		createIncomeHandler := htmxcreateincomehandler.New(factory)

		err := createIncomeHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		errorValue := fmt.Sprintf("no existe una cuenta con el id %s", accountID)
		assert.True(t, strings.Contains(responseBody, errorValue))
		repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("be able to create income", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetIncomeRepositoryMock()
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)

		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepository)

		accountRepository := factory.GetAccountRepositoryMock()
		accountRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepository)
		accountRepository.EXPECT().GetByID(mock.Anything, account.ID()).Return(account, nil)
		accountRepository.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		categoryRepository := factory.GetCategoryRepositoryMock()
		categoryRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithUserID(user.ID()).WithIncomeType().Create(categoryRepository)
		categoryRepository.EXPECT().GetByID(mock.Anything, category.ID()).Return(category, nil)

		ctxBuilder := builders.NewFiberContextBuilder().
			WithMethod("POST").
			WithContentType(fiber.MIMEApplicationForm).
			WithBody(fmt.Appendf(nil,
				"amount=%s&date=%s&category_id=%s",
				"1000",
				account.CreatedAt().Add(time.Hour*25).Format("2006-01-02"),
				category.ID(),
			)).
			WithUser(user)
		defer ctxBuilder.Release()
		ctx := ctxBuilder.Build()

		patches := gomonkey.ApplyMethodReturn(ctx, "Params", account.ID())
		defer patches.Reset()
		createIncomeHandler := htmxcreateincomehandler.New(factory)

		err := createIncomeHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		assert.Nil(t, err)
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Monto: <span>%s</span></p>", "1000")))
	})
}
