package accounthandlers

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
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxcreateaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxgetaccounthandler"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/htmxgetaccountshandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/mocks"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestCreateAccountHTMXHandlerShould(t *testing.T) {
	t.Run("return error when initial balance is not valid", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		initialBalance := "some value"
		ctxBuilder.WithBody(fmt.Appendf(nil,
			"name=%s&initial_balance=%s",
			"test",
			initialBalance,
		))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())
		errorValue := fmt.Sprintf("%s is not valid money value", initialBalance)
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("error creating account: %s", errorValue), odinError.ExternalError())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, errorValue))
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
	})

	t.Run("return error when name is empty", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody(fmt.Appendf(nil,
			"name=%s&initial_balance=%s",
			"",
			"1000000",
		))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		errorValue := "validation error: name cannot be empty"
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("error creating account: %s", errorValue), odinError.ExternalError())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, errorValue))
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
	})

	t.Run("return error when render fails on success", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody(fmt.Appendf(nil,
			"name=%s&initial_balance=%s",
			"test",
			"1000000",
		))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()
		renderError := errors.New("some render error")
		patches := gomonkey.ApplyMethodReturn(ctx, "Render", renderError)
		defer patches.Reset()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("error rendering create account block: %s", renderError), odinError.Error())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.False(t, strings.Contains(responseBody, renderError.Error()))
		assert.Equal(t, odinerrors.RENDER, odinError.Tag())
	})

	t.Run("return error when render fails on error", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		ctxBuilder.WithBody(fmt.Appendf(nil,
			"name=%s&initial_balance=%s",
			"test",
			"some value",
		))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()
		renderError := errors.New("some render error")
		patches := gomonkey.ApplyMethodReturn(ctx, "Render", renderError)
		defer patches.Reset()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())

		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("error rendering create account error block: %s", renderError), odinError.Error())
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.False(t, strings.Contains(responseBody, renderError.Error()))
		assert.Equal(t, odinerrors.RENDER, odinError.Tag())
	})

	t.Run("be able to create an account", func(t *testing.T) {
		repository := mocks.NewMockAccountRepository(t)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("POST").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		accountName := "test"
		initialBalance := "1000000"
		ctxBuilder.WithBody(fmt.Appendf(nil,
			"name=%s&initial_balance=%s",
			accountName,
			initialBalance,
		))
		createAccountHandler := htmxcreateaccounthandler.New(repository)
		ctx := ctxBuilder.Build()

		err := createAccountHandler.Handle(ctx)
		responseBody := string(ctx.Response().Body())

		assert.Nil(t, err)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, string(ctx.Response().Header.ContentType()))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", accountName)))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", initialBalance)))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", initialBalance)))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Fecha apertura: <span>%s</span></p>", time.Now().Format("Monday, _2 January 2006"))))
	})
}

func TestGetAccountsHTMXHandlerShould(t *testing.T) {

	t.Run("does not return any accounts when user has not yet created one", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetAccountRepositoryMock()
		repository.EXPECT().GetAll(mock.Anything).Return([]*accountmodel.Account{}, nil)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		getAccountsHandler := htmxgetaccountshandler.New(repository)
		ctx := ctxBuilder.Build()

		err := getAccountsHandler.Handle(ctx)

		assert.Nil(t, err)
		repository.AssertCalled(t, "GetAll", mock.Anything)
	})

	t.Run("return error when repository returns error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetAccountRepositoryMock()
		repositoryError := errors.New("some repository error")
		repository.EXPECT().GetAll(mock.Anything).Return([]*accountmodel.Account{}, repositoryError)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		getAccountsHandler := htmxgetaccountshandler.New(repository)
		ctx := ctxBuilder.Build()

		err := getAccountsHandler.Handle(ctx)

		assert.Equal(t, repositoryError.Error(), err.Error())
		repository.AssertCalled(t, "GetAll", mock.Anything)
	})

	t.Run("return accounts", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetAccountRepositoryMock()
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		account0 := builders.NewAccountBuilder().WithName("saving account").WithUserID(ctxBuilder.User().ID()).Build()
		account1 := builders.NewAccountBuilder().WithName("cash").WithUserID(ctxBuilder.User().ID()).Build()
		accounts := []*accountmodel.Account{account0, account1}
		repository.EXPECT().GetAll(mock.Anything).Return(accounts, nil)
		getAccountsHandler := htmxgetaccountshandler.New(repository)
		ctx := ctxBuilder.Build()

		err := getAccountsHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		assert.Nil(t, err)
		repository.AssertCalled(t, "GetAll", mock.Anything)
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", account0.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", account0.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", account0.Balance()))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Fecha apertura: <span>%s</span></p>", account0.CreatedAt().Format("Monday, _2 January 2006"))))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", account1.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", account1.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", account1.Balance()))
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Fecha apertura: <span>%s</span></p>", account1.CreatedAt().Format("Monday, _2 January 2006"))))
	})

	t.Run("return error when render fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetAccountRepositoryMock()
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		account0 := builders.NewAccountBuilder().WithName("saving account").WithUserID(ctxBuilder.User().ID()).Build()
		account1 := builders.NewAccountBuilder().WithName("cash").WithUserID(ctxBuilder.User().ID()).Build()
		accounts := []*accountmodel.Account{account0, account1}
		repository.EXPECT().GetAll(mock.Anything).Return(accounts, nil)
		getAccountsHandler := htmxgetaccountshandler.New(repository)
		ctx := ctxBuilder.Build()
		renderError := errors.New("some render error")
		patches := gomonkey.ApplyMethodReturn(ctx, "Render", renderError)
		defer patches.Reset()

		err := getAccountsHandler.Handle(ctx)

		assert.Equal(t, renderError, err)
	})
}

func TestGetAccountHTMXHandlerShould(t *testing.T) {

	t.Run("return error when repository returns error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetAccountRepositoryMock()
		repositoryError := errors.New("some repository error")
		repository.EXPECT().GetByID(mock.Anything, mock.Anything).Return(nil, repositoryError)
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		getAccountHandler := htmxgetaccounthandler.New(repository)
		ctx := ctxBuilder.Build()
		patches := gomonkey.ApplyMethodReturn(ctx, "Params", "some id")
		defer patches.Reset()

		err := getAccountHandler.Handle(ctx)

		assert.Equal(t, repositoryError.Error(), err.Error())
		repository.AssertCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("return error when account does not exist", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepository)

		accountRepository := factory.GetAccountRepositoryMock()
		accountRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepository)
		accountRepository.EXPECT().GetByID(mock.Anything, "some id").Return(nil, errors.New("account not found"))
		ctxBuilder := builders.NewFiberContextBuilder()
		ctxBuilder.WithMethod("GET").WithContentType(fiber.MIMEApplicationForm)
		defer ctxBuilder.Release()
		getAccountHandler := htmxgetaccounthandler.New(accountRepository)
		ctx := ctxBuilder.Build()
		patches := gomonkey.ApplyMethodReturn(ctx, "Params", "some id")
		defer patches.Reset()

		err := getAccountHandler.Handle(ctx)

		assert.Equal(t, "account not found", err.Error())
		accountRepository.AssertCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("return account when belongs to user", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(userRepository)

		accountRepository := factory.GetAccountRepositoryMock()
		accountRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		account := builders.NewAccountBuilder().WithUserID(user.ID()).Create(accountRepository)
		accountRepository.EXPECT().GetByID(mock.Anything, account.ID()).Return(account, nil)
		ctxBuilder := builders.NewFiberContextBuilder().
			WithMethod("GET").
			WithContentType(fiber.MIMEApplicationForm).
			WithUser(user)
		defer ctxBuilder.Release()
		getAccountHandler := htmxgetaccounthandler.New(accountRepository)
		ctx := ctxBuilder.Build()
		patches := gomonkey.ApplyMethodReturn(ctx, "Params", account.ID())
		defer patches.Reset()

		err := getAccountHandler.Handle(ctx)

		responseBody := string(ctx.Response().Body())
		assert.Nil(t, err)
		assert.True(t, strings.Contains(responseBody, fmt.Sprintf("<p>Nombre: <span>%s</span></p>", account.Name())))
	})
}
