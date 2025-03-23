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
