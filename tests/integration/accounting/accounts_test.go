package accounting_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/accountsrepositoryfactory"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

const accountPath = "/accounts"

func TestCreateAccountHtmxShould(t *testing.T) {

	t.Run("create account when everything is ok", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body)

		response, _ := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return unauthorized when request is anonymous", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithAnonymousSession()

		response, _ := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return unauthorized when request sent cookie but session doesn't exists", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		user := userbuilder.New().Create(accountsFactory.GetUserRepository())
		session := sessionmodel.New(user.ID())
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithSession(session)
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)

		response, _ := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})
}

func TestGetAccountsHTMXShould(t *testing.T) {

	t.Run("return accounts when everything is ok", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		user := userbuilder.New().Create(accountsFactory.GetUserRepository())
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithMethod("GET").
			WithPath(accountPath).
			WithContentType("").
			WithUser(user)

		account0 := builders.NewAccountBuilder().
			WithName("saving account").
			WithUserID(user.ID()).
			Create(accountingFactory.GetAccountRepository())
		account1 := builders.NewAccountBuilder().
			WithName("cash").
			WithUserID(user.ID()).
			Create(accountingFactory.GetAccountRepository())
		user1 := userbuilder.New().WithEmail("some@email.com").Create(accountsFactory.GetUserRepository())
		account2 := builders.NewAccountBuilder().
			WithName("nu").
			WithUserID(user1.ID()).
			WithInitialBalance("0").
			Create(accountingFactory.GetAccountRepository())

		response, responseBody := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", account0.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", account0.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", account0.Balance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Fecha apertura: <span>%s</span></p>", account0.CreatedAt().Format("Monday, _2 January 2006")))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", account1.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", account1.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", account1.Balance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<p>Fecha apertura: <span>%s</span></p>", account1.CreatedAt().Format("Monday, _2 January 2006")))
		assert.NotContains(t, responseBody, fmt.Sprintf("<p>Name: <span>%s</span></p>", account2.Name()))
		assert.NotContains(t, responseBody, fmt.Sprintf("<p>Saldo inicial: <span>%s</span></p>", account2.InitialBalance()))
		assert.NotContains(t, responseBody, fmt.Sprintf("<p>Saldo actual: <span>%s</span></p>", account2.Balance()))
	})

	t.Run("return unauthorized when request is anonymous", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithMethod("GET").
			WithPath(accountPath).
			WithContentType("").
			WithAnonymousSession()

		response, _ := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})

	t.Run("return unauthorized when request sent cookie but session doesn't exists", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		user := userbuilder.New().Create(accountsFactory.GetUserRepository())
		session := sessionmodel.New(user.ID())
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithMethod("GET").
			WithPath(accountPath).
			WithContentType("").
			WithSession(session)

		response, _ := testutils.GetHtmlResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})
}
