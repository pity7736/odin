package accounts_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
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

func TestRegisterIntegrationShould(t *testing.T) {
	t.Run("register a new user and then log in with the same credentials", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		email := "julian@example.com"
		authHash := "client-derived-auth-hash"
		registerBody := fmt.Sprintf(
			`{"email": "%s", "auth_hash": "%s", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
			email,
			authHash,
		)
		var registerResponseData map[string]any
		registerRequestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/register").
			WithPayload(registerBody).
			WithResponseData(&registerResponseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		registerResponse := testutils.GetJSONResponseFromRequestBuilder(application, registerRequestBuilder)
		defer func() { _ = registerResponse.Body.Close() }()
		assert.Equal(t, http.StatusCreated, registerResponse.StatusCode)
		assert.NotEmpty(t, registerResponseData["id"])
		assert.Equal(t, email, registerResponseData["email"])
		assert.NotContains(t, registerResponseData, "token")
		assert.NotContains(t, registerResponseData, "encrypted_master_key")
		loginBody := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, email, authHash)
		var loginResponseData map[string]any
		loginRequestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/login").
			WithPayload(loginBody).
			WithResponseData(&loginResponseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		loginResponse := testutils.GetJSONResponseFromRequestBuilder(application, loginRequestBuilder)
		defer func() { _ = loginResponse.Body.Close() }()
		assert.Equal(t, http.StatusCreated, loginResponse.StatusCode)
		assert.NotEmpty(t, loginResponseData["token"])
		assert.Equal(t, "encrypted-master-key-base64", loginResponseData["encrypted_master_key"])
		keyParams, ok := loginResponseData["key_params"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "argon2id", keyParams["algorithm"])
		assert.Equal(t, "salt-base64", keyParams["salt"])
	})
	t.Run("reject registration when the email is already registered", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		email := "julian@example.com"
		userbuilder.New().WithEmail(email).Create(userRepository)
		body := fmt.Sprintf(
			`{"email": "%s", "auth_hash": "client-derived-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
			email,
		)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/register").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusConflict, response.StatusCode)
		assert.Equal(t, "El correo ya está registrado", responseData["error"])
	})
}

func TestLoginIntegrationShould(t *testing.T) {
	t.Run("return token, encrypted master key and key params on successful login", func(t *testing.T) {
		application, userRepository, sessionRepository := newIntegrationApp()
		user := userbuilder.New().WithEmail("some@email.com").Create(userRepository)
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, user.Email(), userbuilder.DefaultPassword)
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
		user := userbuilder.New().WithEmail("some@email.com").Create(userRepository)
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "wrong-password"}`, user.Email())
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
