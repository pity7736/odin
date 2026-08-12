package login_api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApplication(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory, factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func TestRestLoginShould(t *testing.T) {
	t.Run("return session token when credentials are valid", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		factory.GetUserRepositoryMock().EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		factory.GetSessionRepositoryMock().EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email(), builder.Password())
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotEmpty(t, responseData["token"])
		assert.Empty(t, responseData["error"])
	})
	t.Run("return error when user does not exist", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetUserRepositoryMock().EXPECT().GetByEmail(mock.Anything, "unknown@example.com").Return(nil, nil)
		body := `{"email": "unknown@example.com", "password": "any"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Empty(t, responseData["token"])
		assert.Equal(t, "Correo o contraseña incorrectos", responseData["error"])
	})
	t.Run("return error when password is wrong", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		factory.GetUserRepositoryMock().EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		body := fmt.Sprintf(`{"email": "%s", "password": "wrong_password"}`, user.Email())
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Correo o contraseña incorrectos", responseData["error"])
	})
	t.Run("return error when email is missing", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		body := `{"password": "some"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "El correo es obligatorio", responseData["error"])
	})
	t.Run("return error when email is empty", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		body := `{"email": "", "password": "some"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "El correo es obligatorio", responseData["error"])
	})
	t.Run("return error when password is missing", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		body := fmt.Sprintf(`{"email": "%s"}`, builder.Build().Email())
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "La contraseña es obligatoria", responseData["error"])
	})
	t.Run("return error when password is empty", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		body := fmt.Sprintf(`{"email": "%s", "password": ""}`, builder.Build().Email())
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "La contraseña es obligatoria", responseData["error"])
	})
	t.Run("return error when body is malformed", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		body := `{"email": "test@example.com" "password": "x"}`
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		var responseData map[string]string
		response := testutils.GetJSONResponseFromRequest(newApplication(factory), req, &responseData)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Datos de solicitud inválidos", responseData["error"])
	})
}

func TestHtmxLoginShould(t *testing.T) {
	t.Run("render login form", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		response, body := testutils.GetHTMLResponseFromRequest(newApplication(factory), req)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Contains(t, body, `<form hx-post="/auth/login?next=/" hx-target="#login_error">`)
	})
	t.Run("set session cookie and redirect when credentials are valid", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		factory.GetUserRepositoryMock().EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		factory.GetSessionRepositoryMock().EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		body := fmt.Sprintf("email=%s&password=%s", user.Email(), builder.Password())
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationForm)
		response, _ := testutils.GetHTMLResponseFromRequest(newApplication(factory), req)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "/", response.Header.Get("HX-Redirect"))
		sessionCookie := response.Cookies()[0]
		assert.NotEmpty(t, sessionCookie.Value)
		assert.True(t, sessionCookie.Secure)
		assert.True(t, sessionCookie.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, sessionCookie.SameSite)
	})
	t.Run("render error when user does not exist", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetUserRepositoryMock().EXPECT().GetByEmail(mock.Anything, "nobody@example.com").Return(nil, nil)
		body := "email=nobody@example.com&password=any"
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationForm)
		response, responseBody := testutils.GetHTMLResponseFromRequest(newApplication(factory), req)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Contains(t, responseBody, "Correo o contraseña incorrectos")
	})
}
