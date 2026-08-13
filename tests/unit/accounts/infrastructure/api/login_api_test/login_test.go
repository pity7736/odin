package login_api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApplication(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory, factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func TestRest(t *testing.T) {
	t.Run("non existing email", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		email := "some@email.com"
		body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, builder.Password())
		var responseData map[string]string
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

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Correo o contraseña incorrectos", responseData["error"])
		assert.Empty(t, responseData["token"])
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
				fmt.Sprintf(`{"password": "%s"}`, builder.Password()),
				"El correo es obligatorio",
			},
			{
				"when email is empty",
				fmt.Sprintf(`{"email": "", "password": "%s"}`, builder.Password()),
				"El correo es obligatorio",
			},
			{
				"when password is missing",
				fmt.Sprintf(`{"email": "%s"}`, user.Email()),
				"La contraseña es obligatoria",
			},
			{
				"when password is empty",
				fmt.Sprintf(`{"email": "%s", "password": ""}`, user.Email()),
				"La contraseña es obligatoria",
			},
			{
				"when body is wrong",
				fmt.Sprintf(`{"email": "%s" "password": ""}`, user.Email()),
				"Datos de solicitud inválidos",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				var responseData map[string]string
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
				assert.Empty(t, responseData["token"])
				repository.AssertNotCalled(t, "GetByEmail")
			})
		}
	})

	t.Run("when email and password are correct", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email(), builder.Password())
		var responseData map[string]string
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
		assert.Empty(t, responseData["error"])
		assert.NotEmpty(t, responseData["token"])
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})
}

func TestHTMX(t *testing.T) {
	t.Run("get login form", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/auth/login").
			WithMethod("GET").
			WithContentType("").
			WithAnonymousSession()
		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.True(t, strings.Contains(responseData, `<p>Login</p>`))
		assert.True(t, strings.Contains(responseData, `<form hx-post="/auth/login?next=/" hx-target="#login_error">`))
		assert.True(t, strings.Contains(responseData, `<label for="email">Email:</label>`))
		assert.True(t, strings.Contains(responseData, `<input id="email" type="email" name="email" required>`))
		assert.True(t, strings.Contains(responseData, `<label for="password">Password:</label>`))
		assert.True(t, strings.Contains(responseData, `<input id="password" type="password" name="password" required>`))
		assert.True(t, strings.Contains(responseData, `<button type="submit">Iniciar sesión</button>`))
		assert.True(t, strings.Contains(responseData, `</form>`))
	})

	t.Run("non existing email", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		email := "some@email.com"
		body := fmt.Sprintf("email=%s&password=%s", email, builder.Password())
		repository := factory.GetUserRepositoryMock()
		repository.EXPECT().GetByEmail(mock.Anything, email).Return(nil, nil)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/auth/login").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationForm).
			WithAnonymousSession()
		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.True(t, strings.Contains(responseData, "Correo o contraseña incorrectos"))
		repository.AssertCalled(t, "GetByEmail", mock.Anything, email)
	})

	t.Run("when email and password are correct", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := newApplication(factory)
		builder := userbuilder.New()
		user := builder.Build()
		body := fmt.Sprintf("email=%s&password=%s", user.Email(), builder.Password())
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		sessionRepositoryMock := factory.GetSessionRepositoryMock()
		sessionRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		requestBuilder := builders.NewRequestBuilder(factory.GetUserRepository(), factory.GetSessionRepository()).
			WithPath("/auth/login").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationForm).
			WithAnonymousSession()
		response, _ := testutils.GetHTMLResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		sessionCookie := response.Cookies()[0]

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotEmpty(t, sessionCookie.Value)
		assert.True(t, sessionCookie.Secure)
		assert.True(t, sessionCookie.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, sessionCookie.SameSite)
		assert.Equal(t, "/", response.Header.Get("HX-Redirect"))
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})
}
