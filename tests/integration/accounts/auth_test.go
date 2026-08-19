package accounts_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/testutils"
)

func newIntegrationApp() (app.Application, *inmemory.InMemoryUserRepository, *inmemory.InMemorySessionRepository) {
	userRepository := inmemory.NewInMemoryUserRepository()
	sessionRepository := inmemory.NewInMemorySessionRepository()
	application := app.NewFiberApplication(
		sessionRepository,
		userRepository,
		bcrypthasher.New(),
	)
	return application, userRepository, sessionRepository
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
