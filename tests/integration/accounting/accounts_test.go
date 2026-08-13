package accounting_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/pgrepositories"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func newIntegrationApp() (app.Application, accountingrepositoryfactory.RepositoryFactory, *pgrepositories.PGUserRepository, *pgrepositories.PGSessionRepository) {
	accountingFactory := accountingrepositoryfactory.New()
	userRepository := pgrepositories.NewPGUserRepository()
	sessionRepository := pgrepositories.NewPGSessionRepository()
	application := app.NewFiberApplication(
		accountingFactory,
		sessionRepository,
		userRepository,
		bcrypthasher.New(),
	)
	return application, accountingFactory, userRepository, sessionRepository
}

const accountPath = "/accounts"

func TestCreateAccountHtmxShould(t *testing.T) {
	t.Run("create account when everything is ok", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		body := fmt.Sprintf("name=%s&initial_balance=%s", "test", "10000")
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body)

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return bad request when data is wrong", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		body := fmt.Sprintf("name=%s&initial_balance=%s", "test", "aoeu")
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body)

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return redirect when request is anonymous", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		body := fmt.Sprintf("name=%s&initial_balance=%s", "test", "10000")
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithAnonymousSession()

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusFound, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return redirect when session does not exist", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		user := userbuilder.New().Create(userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		body := fmt.Sprintf("name=%s&initial_balance=%s", "test", "10000")
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithSession(session)

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusFound, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})
}

func TestGetAccountsHTMXShould(t *testing.T) {
	t.Run("return accounts when everything is ok", func(t *testing.T) {
		application, accountingFactory, userRepository, sessionRepository := newIntegrationApp()
		user := userbuilder.New().Create(userRepository)
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
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
		user1 := userbuilder.New().WithEmail("some@email.com").Create(userRepository)
		account2 := builders.NewAccountBuilder().
			WithName("nu").
			WithUserID(user1.ID()).
			WithInitialBalance("0").
			Create(accountingFactory.GetAccountRepository())

		response, responseBody := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		assert.Contains(t, responseBody, fmt.Sprintf("<td><a href=\"/accounts/%s\">%s</a></td>", account0.ID(), account0.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account0.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account0.Balance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account0.CreatedAt().Format("Monday, _2 January 2006")))

		assert.Contains(t, responseBody, fmt.Sprintf("<td><a href=\"/accounts/%s\">%s</a></td>", account1.ID(), account1.Name()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account1.InitialBalance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account1.Balance()))
		assert.Contains(t, responseBody, fmt.Sprintf("<td>%s</td>", account1.CreatedAt().Format("Monday, _2 January 2006")))

		assert.NotContains(t, responseBody, fmt.Sprintf("<td><a href=\"/accounts/%s\">%s<a></td>", account2.ID(), account2.Name()))
		assert.NotContains(t, responseBody, fmt.Sprintf("<td>%s</td>", account2.InitialBalance()))
		assert.NotContains(t, responseBody, fmt.Sprintf("<td>%s</td>", account2.Balance()))
	})

	t.Run("return found when request is anonymous", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath(accountPath).
			WithContentType("").
			WithAnonymousSession()

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusFound, response.StatusCode)
	})

	t.Run("return found when session does not exist", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		user := userbuilder.New().Create(userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath(accountPath).
			WithContentType("").
			WithSession(session)

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, fiber.StatusFound, response.StatusCode)
	})

	t.Run("redirect to login when not authenticated", func(t *testing.T) {
		application, _, _, _ := newIntegrationApp()
		req := httptest.NewRequest("GET", "/accounts", nil)
		response, err := application.Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
	})
}

func TestAccountIntegrationShould(t *testing.T) {
	t.Run("create an account when authenticated", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository)
		body := `{"name": "Ahorros", "initial_balance": "100000"}`
		requestBuilder.
			WithMethod("POST").
			WithPath("/api/v1/accounts").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("return error when account name is empty", func(t *testing.T) {
		application, _, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository)
		body := `{"name": "", "initial_balance": "100000"}`
		requestBuilder.
			WithMethod("POST").
			WithPath("/api/v1/accounts").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
