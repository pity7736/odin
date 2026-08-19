package login_api_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApplication(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func TestRest(t *testing.T) {
	t.Run("non existing email", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		email := "some@email.com"
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, email, builder.Password())
		var responseData map[string]any
		repository := factory.GetUserRepositoryMock()
		repository.EXPECT().GetByEmail(mock.Anything, email).Return(nil, nil)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
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
		repository.AssertCalled(t, "GetByEmail", mock.Anything, email)
	})

	t.Run("login with wrong data", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		testCases := []struct {
			name          string
			body          string
			expectedError string
		}{
			{
				"when email is missing",
				fmt.Sprintf(`{"auth_hash": "%s"}`, builder.Password()),
				"El correo es obligatorio",
			},
			{
				"when email is empty",
				fmt.Sprintf(`{"email": "", "auth_hash": "%s"}`, builder.Password()),
				"El correo es obligatorio",
			},
			{
				"when auth hash is missing",
				fmt.Sprintf(`{"email": "%s"}`, user.Email()),
				"La contraseña es obligatoria",
			},
			{
				"when auth hash is empty",
				fmt.Sprintf(`{"email": "%s", "auth_hash": ""}`, user.Email()),
				"La contraseña es obligatoria",
			},
			{
				"when body is wrong",
				fmt.Sprintf(`{"email": "%s" "auth_hash": ""}`, user.Email()),
				"Datos de solicitud inválidos",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				var responseData map[string]any
				repository := factory.GetUserRepositoryMock()
				requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
					WithPath("/api/v1/auth/login").
					WithPayload(testCase.body).
					WithResponseData(&responseData).
					WithContentType(fiber.MIMEApplicationJSON).
					WithAnonymousSession()

				response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
				defer func() { _ = response.Body.Close() }()

				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Equal(t, testCase.expectedError, responseData["error"])
				assert.NotContains(t, responseData, "token")
				assert.NotContains(t, responseData, "encrypted_master_key")
				assert.NotContains(t, responseData, "key_params")
				repository.AssertNotCalled(t, "GetByEmail")
			})
		}
	})

	t.Run("when email and auth hash are correct", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, user.Email(), builder.Password())
		var responseData map[string]any
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		sessionRepositoryMock := factory.GetSessionRepositoryMock()
		sessionRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotContains(t, responseData, "error")
		assert.NotEmpty(t, responseData["token"])
		assert.Equal(t, userbuilder.DefaultEncryptedMasterKey, responseData["encrypted_master_key"])
		keyParams, ok := responseData["key_params"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, userbuilder.DefaultAlgorithm, keyParams["algorithm"])
		assert.Equal(t, float64(userbuilder.DefaultIterations), keyParams["iterations"])
		assert.Equal(t, float64(userbuilder.DefaultMemory), keyParams["memory"])
		assert.Equal(t, float64(userbuilder.DefaultParallelism), keyParams["parallelism"])
		assert.Equal(t, userbuilder.DefaultSalt, keyParams["salt"])
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})

	t.Run("when the use case returns a raw non-odin error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, user.Email(), builder.Password())
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(nil, errors.New("database unavailable"))
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Zero(t, response.ContentLength)
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})

	t.Run("when the use case returns an odin error with a non-4xx handled tag", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, user.Email(), builder.Password())
		var responseData map[string]any
		renderError := odinerrors.NewErrorBuilder("render failed").
			WithExternalMessage("Error interno").
			WithTag(odinerrors.Render).
			Build()
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(nil, renderError)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Equal(t, "Error interno", responseData["error"])
		assert.NotContains(t, responseData, "token")
		assert.NotContains(t, responseData, "encrypted_master_key")
		assert.NotContains(t, responseData, "key_params")
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})

	t.Run("when the use case returns an odin error with a domain tag", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf(`{"email": "%s", "auth_hash": "%s"}`, user.Email(), builder.Password())
		var responseData map[string]any
		domainError := odinerrors.NewErrorBuilder("domain failure").
			WithExternalMessage("Solicitud inválida").
			WithTag(odinerrors.Domain).
			Build()
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(nil, domainError)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/api/v1/auth/login").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Solicitud inválida", responseData["error"])
		assert.NotContains(t, responseData, "token")
		assert.NotContains(t, responseData, "encrypted_master_key")
		assert.NotContains(t, responseData, "key_params")
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})
}
