package accounting_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/pgrepositories"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/testutils"
)

func newIntegrationApp() (app.Application, *pgrepositories.PGUserRepository, *pgrepositories.PGSessionRepository) {
	userRepository := pgrepositories.NewPGUserRepository()
	sessionRepository := pgrepositories.NewPGSessionRepository()
	application := app.NewFiberApplication(
		accountingrepositoryfactory.New(),
		sessionRepository,
		userRepository,
	)
	return application, userRepository, sessionRepository
}

func TestAccountIntegrationShould(t *testing.T) {
	t.Run("create an account when authenticated", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
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
		application, userRepository, sessionRepository := newIntegrationApp()
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

func TestLogoutIntegrationShould(t *testing.T) {
	t.Run("terminate session when authenticated via REST", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository)
		requestBuilder.
			WithMethod("DELETE").
			WithPath("/api/v1/auth/logout").
			WithContentType(fiber.MIMEApplicationJSON)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

func TestCategoryIntegrationShould(t *testing.T) {
	t.Run("create a category when authenticated", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository)
		body := `{"name": "Comida", "type": "expense"}`
		requestBuilder.
			WithMethod("POST").
			WithPath("/api/v1/categories").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("retrieve categories when authenticated", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository)
		requestBuilder.
			WithMethod("GET").
			WithPath("/api/v1/categories").
			WithContentType(fiber.MIMEApplicationJSON)
		var responseData map[string]any
		requestBuilder.WithResponseData(&responseData)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("redirect to login when not authenticated", func(t *testing.T) {
		application, _, _ := newIntegrationApp()
		req := httptest.NewRequest("GET", "/accounts", nil)
		response, err := application.Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
	})
}
