package register_api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/apptest"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApplication(factory *testrepositoryfactory.Factory) app.Application {
	return apptest.New().
		WithSessionRepository(factory.GetSessionRepository()).
		WithUserRepository(factory.GetUserRepository()).
		Build()
}

func validRegisterBody(email string) string {
	return fmt.Sprintf(
		`{"email": "%s", "auth_hash": "client-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
		email,
	)
}

func TestRegisterRestShould(t *testing.T) {
	t.Run("create a new account", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		email := "new@example.com"
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().Exists(mock.Anything, email).Return(false, nil)
		userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/users").
			WithPayload(validRegisterBody(email)).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotEmpty(t, responseData["id"])
		assert.Equal(t, email, responseData["email"])
		assert.NotContains(t, responseData, "token")
		assert.NotContains(t, responseData, "encrypted_master_key")
		userRepository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject registration when the email already exists", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		existingUser := userbuilder.New().WithEmail("taken@example.com").Build()
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().Exists(mock.Anything, existingUser.Email()).Return(true, nil)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/users").
			WithPayload(validRegisterBody(existingUser.Email())).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusConflict, response.StatusCode)
		assert.Equal(t, "El correo ya está registrado", responseData["error"])
		userRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject invalid registration data", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		testCases := []struct {
			name          string
			body          string
			expectedError string
		}{
			{
				"when email is empty",
				`{"email": "", "auth_hash": "client-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
				"El correo es obligatorio",
			},
			{
				"when auth hash is empty",
				`{"email": "new@example.com", "auth_hash": "", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
				"La contraseña es obligatoria",
			},
			{
				"when encrypted master key is empty",
				`{"email": "new@example.com", "auth_hash": "client-auth-hash", "encrypted_master_key": "", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
				"Datos de solicitud inválidos",
			},
			{
				"when key params algorithm is empty",
				`{"email": "new@example.com", "auth_hash": "client-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "", "salt": "salt-base64", "iterations": 3, "memory": 65536, "parallelism": 4}}`,
				"El algoritmo es obligatorio",
			},
			{
				"when key params iterations is zero",
				`{"email": "new@example.com", "auth_hash": "client-auth-hash", "encrypted_master_key": "encrypted-master-key-base64", "key_params": {"algorithm": "argon2id", "salt": "salt-base64", "iterations": 0, "memory": 65536, "parallelism": 4}}`,
				"Las iteraciones deben ser positivas",
			},
			{
				"when body is malformed",
				`{"email": "new@example.com" "auth_hash": "client-auth-hash"}`,
				"Datos de solicitud inválidos",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				var responseData map[string]any
				userRepository := factory.GetUserRepositoryMock()
				requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
					WithPath("/api/v1/users").
					WithPayload(testCase.body).
					WithResponseData(&responseData).
					WithContentType(fiber.MIMEApplicationJSON).
					WithAnonymousSession()
				response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
				defer func() { _ = response.Body.Close() }()
				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Equal(t, testCase.expectedError, responseData["error"])
				userRepository.AssertNotCalled(t, "Exists")
				userRepository.AssertNotCalled(t, "Add")
			})
		}
	})
}
