package accounthandlers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxcreateaccounthandler"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/unit/mocks"
)

func TestCreateAccountHTMXHandlerShould(t *testing.T) {
	t.Run("return error when initial balance is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		initialBalance := "some value"
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			initialBalance,
		)))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())
		errorValue := fmt.Sprintf("%s is not valid money value", initialBalance)

		assert.Equal(t, errorValue, err.Error())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, errorValue))
	})

	t.Run("return error when name is empty", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"",
			"1000000",
		)))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())
		errorValue := "name cannot be empty"

		assert.Equal(t, errorValue, err.Error())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, errorValue))
	})

	t.Run("return error when render fails", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"1000000",
		)))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()
		renderError := errors.New("some render error")
		patches := gomonkey.ApplyMethodReturn(ctx, "Render", renderError)
		defer patches.Reset()

		err := createAccountHandler.Handle(ctx)

		assert.Equal(t, renderError, err)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
	})

	t.Run("be able to create an account", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		accountName := "test"
		initialBalance := "1000000"
		ctxBuilder.WithBody([]byte(fmt.Sprintf(
			"name=%s&initial_balance=%s",
			accountName,
			initialBalance,
		)))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())

		assert.Nil(t, err)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", accountName)))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Initial Balance: <span>%s</span></p>", initialBalance)))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Balance: <span>%s</span></p>", initialBalance)))
	})
}
