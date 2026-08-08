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
			`{"name":"%s","initial_balance":"%s"}`,
			account.Name(),
			account.InitialBalance().Value(),
		)))
		repository.EXPECT().Add(mock.Anything, account).Return(nil)
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var responseBody map[string]string
		json.Unmarshal(ctx.Response().Body(), &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
		assert.Equal(t, account.Name(), responseBody["name"])
		assert.Equal(t, account.InitialBalance().String(), responseBody["initial_balance"])
		assert.Equal(t, account.Balance().String(), responseBody["balance"])
		assert.Equal(t, account.ID(), responseBody["id"])
		assert.Equal(t, account.UserID(), responseBody["user_id"])
		assert.Equal(t, account.CreatedAt().Format(time.RFC3339), responseBody["created_at"])
		repository.AssertCalled(t, "Add", mock.Anything, account)
	})

	t.Run("not overwrite buffer when creating multiple accounts", func(t *testing.T) {
		// 1. Setup
		repository := mocks.NewMockAccountRepository(t)
		createAccountHandler := restcreateaccounthandler.New(repository)

		// This slice will store pointers to the accounts passed to repository.Add
		var capturedAccounts []*accountmodel.Account

		// We configure the mock to capture the account object from each call
		repository.On("Add", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			account := args.Get(1).(*accountmodel.Account)
			capturedAccounts = append(capturedAccounts, account)
		}).Return(nil)

		// 2. First request to create "nequi"
		ctxBuilder1 := builders.NewFiberContextBuilder()
		ctxBuilder1.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		ctxBuilder1.WithBody([]byte(`{"name": "nequi", "initial_balance": "100"}`))
		defer ctxBuilder1.Release()

		err1 := createAccountHandler.Handle(ctxBuilder1.Build())
		assert.NoError(t, err1)

		// 3. Second request to create "nu"
		ctxBuilder2 := builders.NewFiberContextBuilder()
		ctxBuilder2.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		ctxBuilder2.WithBody([]byte(`{"name": "nu", "initial_balance": "200"}`))
		defer ctxBuilder2.Release()

		err2 := createAccountHandler.Handle(ctxBuilder2.Build())
		assert.NoError(t, err2)

		// 4. Verification
		// Check that two accounts were "saved"
		assert.Equal(t, 2, len(capturedAccounts), "Expected two accounts to be captured")

		// The crucial check: The first captured account's name must NOT have changed.
		// If the fix is not applied, this will fail with:
		// "Expected: 'nequi', Actual: 'nuqui'"
		assert.Equal(t, "nequi", capturedAccounts[0].Name(), "The first account's name was overwritten")
		assert.Equal(t, "nu", capturedAccounts[1].Name(), "The second account's name is incorrect")

		repository.AssertNumberOfCalls(t, "Add", 2)
	})

	t.Run("return error when initial balance is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationJSON)
		defer ctxBuilder.Release()
		initialBalance := "some value"
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			`{"name":"%s","initial_balance":"%s"}`,
			"test",
			initialBalance,
		)))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf(`%s is not valid money value`, initialBalance), odinError.ExternalError())
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
	})

	t.Run("return error when body is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType("application/json")
		defer ctxBuilder.Release()
		initialBalance := "some value"
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			`"name":"%s","initial_balance": %s"`,
			"test",
			initialBalance,
		)))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		assert.NotNil(t, err)
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
			`{"name":"%s","initial_balance":"%s"}`,
			account.Name(),
			account.InitialBalance().Value(),
		)))
		repository.EXPECT().Add(mock.Anything, account).Return(errors.New("some error"))
		createAccountHandler := restcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, fiber.MIMEApplicationJSON, string(ctx.Response().Header.ContentType()))
	})

}
