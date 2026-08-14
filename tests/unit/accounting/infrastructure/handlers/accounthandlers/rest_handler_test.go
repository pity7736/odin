package accounthandlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/restcreateaccounthandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/unit/mocks"
)

func TestCreateAccountHandlerShould(t *testing.T) {
	t.Run("be able to create an account", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType("application/json")
		defer ctxBuilder.Release()
		id, _ := uuid.NewV7()
		patches := gomonkey.ApplyFuncReturn(uuid.NewV7, id, nil)
		defer patches.Reset()
		account := builders.NewAccountBuilder().WithUserID(ctxBuilder.User().ID()).Build()
		timePatches := gomonkey.ApplyFuncReturn(time.Now, account.CreatedAt())
		defer timePatches.Reset()
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			`{"name":"%s","initial_balance":"%s","type":"savings","currency":"COP"}`,
			account.Name(),
			account.InitialBalance().Value(),
		)))
		repository.EXPECT().ExistsByNameAndCurrency(mock.Anything, account.Name(), account.Currency()).Return(false, nil)
		repository.EXPECT().Add(mock.Anything, account).Return(nil)
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var responseBody map[string]string
		_ = json.Unmarshal(ctx.Response().Body(), &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
		assert.Equal(t, account.Name(), responseBody["name"])
		assert.Equal(t, account.InitialBalance().String(), responseBody["initial_balance"])
		assert.Equal(t, account.Balance().String(), responseBody["balance"])
		assert.Equal(t, account.Type().String(), responseBody["type"])
		assert.Equal(t, account.Currency().String(), responseBody["currency"])
		assert.Equal(t, account.ID(), responseBody["id"])
		assert.Equal(t, account.UserID(), responseBody["user_id"])
		assert.Equal(t, account.CreatedAt().Format(time.RFC3339), responseBody["created_at"])
		repository.AssertCalled(t, "Add", mock.Anything, account)
	})

	t.Run("not overwrite buffer when creating multiple accounts", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		createAccountHandler := restcreateaccounthandler.New(repository)

		var capturedAccounts []*accountmodel.Account

		repository.On("ExistsByNameAndCurrency", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
		repository.On("Add", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			account := args.Get(1).(*accountmodel.Account)
			capturedAccounts = append(capturedAccounts, account)
		}).Return(nil)

		ctxBuilder1 := builders.NewFiberContextBuilder()
		ctxBuilder1.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		ctxBuilder1.WithBody([]byte(`{"name": "nequi", "initial_balance": "100", "type": "savings", "currency": "COP"}`))
		defer ctxBuilder1.Release()

		err1 := createAccountHandler.Handle(ctxBuilder1.Build())
		assert.NoError(t, err1)

		ctxBuilder2 := builders.NewFiberContextBuilder()
		ctxBuilder2.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		ctxBuilder2.WithBody([]byte(`{"name": "nu", "initial_balance": "200", "type": "savings", "currency": "COP"}`))
		defer ctxBuilder2.Release()

		err2 := createAccountHandler.Handle(ctxBuilder2.Build())
		assert.NoError(t, err2)

		assert.Equal(t, 2, len(capturedAccounts), "Expected two accounts to be captured")
		assert.Equal(t, "nequi", capturedAccounts[0].Name(), "The first account's name was overwritten")
		assert.Equal(t, "nu", capturedAccounts[1].Name(), "The second account's name is incorrect")
		repository.AssertNumberOfCalls(t, "Add", 2)
	})

	t.Run("return error when name is empty", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`{"name":"","initial_balance":"1000000","type":"savings","currency":"COP"}`))
		repository.EXPECT().ExistsByNameAndCurrency(mock.Anything, "", mock.Anything).Return(false, nil)
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El nombre es obligatorio", odinError.ExternalError())
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when initial balance is missing", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`{"name":"test","type":"savings","currency":"COP"}`))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El saldo inicial es obligatorio", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when initial balance is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`{"name":"test","initial_balance":"some value","type":"savings","currency":"COP"}`))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, "El valor ingresado no es un monto válido", odinError.ExternalError())
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
	})

	t.Run("return error when type is invalid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`{"name":"test","initial_balance":"100","type":"invalid","currency":"COP"}`))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Tipo de cuenta inválido", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when currency is invalid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`{"name":"test","initial_balance":"100","type":"savings","currency":"EUR"}`))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Moneda inválida", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when body is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType("application/json")
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(`"name":"test","initial_balance": some value"`))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		assert.Contains(t, err.Error(), "invalid character")
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
	})

	t.Run("return error when repository returns error", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType("application/json")
		defer ctxBuilder.Release()
		id, _ := uuid.NewV7()
		patches := gomonkey.ApplyFuncReturn(uuid.NewV7, id, nil)
		defer patches.Reset()
		account := builders.NewAccountBuilder().WithUserID(ctxBuilder.User().ID()).Build()
		timePatches := gomonkey.ApplyFuncReturn(time.Now, account.CreatedAt())
		defer timePatches.Reset()
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			`{"name":"%s","initial_balance":"%s","type":"savings","currency":"COP"}`,
			account.Name(),
			account.InitialBalance().Value(),
		)))
		repository.EXPECT().ExistsByNameAndCurrency(mock.Anything, account.Name(), account.Currency()).Return(false, nil)
		repository.EXPECT().Add(mock.Anything, account).Return(errors.New("some error"))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
	})
}
