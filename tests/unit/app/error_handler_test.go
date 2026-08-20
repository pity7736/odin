package app_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/apptest"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func newErrorHandlerApp() (app.Application, *inmemory.InMemoryUserRepository, *inmemory.InMemorySessionRepository) {
	userRepository := inmemory.NewInMemoryUserRepository()
	sessionRepository := inmemory.NewInMemorySessionRepository()
	application := apptest.New().
		WithUserRepository(userRepository).
		WithSessionRepository(sessionRepository).
		Build()
	return application, userRepository, sessionRepository
}

func TestErrorHandlerShould(t *testing.T) {
	t.Run("map an AlreadyExists error to 409 with the external message", func(t *testing.T) {
		application, userRepository, sessionRepository := newErrorHandlerApp()
		email := "taken@example.com"
		userbuilder.New().WithEmail(email).Create(userRepository)
		body := fmt.Sprintf(
			`{"email": "%s", "auth_hash": "client-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
			email,
		)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/users").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusConflict, response.StatusCode)
		assert.Equal(t, "El correo ya está registrado", responseData["error"])
	})
	t.Run("map a Domain error to 400 with the external message", func(t *testing.T) {
		application, userRepository, sessionRepository := newErrorHandlerApp()
		body := `{"email": "", "auth_hash": "client-auth-hash"}`
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "El correo es obligatorio", responseData["error"])
	})
}
