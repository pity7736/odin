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

func TestLoginIntegrationShould(t *testing.T) {
	t.Run("return token, encrypted master key and key params on successful login", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		body := `{"email": "some@email.com", "auth_hash": "password"}`
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotEmpty(t, responseData["token"])
		assert.NotEmpty(t, responseData["encrypted_master_key"])
		keyParams, ok := responseData["key_params"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "argon2id", keyParams["algorithm"])
		assert.NotEmpty(t, keyParams["salt"])
	})
	t.Run("return error on wrong credentials", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		body := `{"email": "some@email.com", "auth_hash": "wrong-password"}`
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, "Correo o contraseña incorrectos", responseData["error"])
		assert.NotContains(t, responseData, "token")
		assert.NotContains(t, responseData, "encrypted_master_key")
		assert.NotContains(t, responseData, "key_params")
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
