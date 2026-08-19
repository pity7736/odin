package logout_api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	handler "raiseexception.dev/odin/src/shared/infrastructure/api"

	"raiseexception.dev/odin/src/accounts/infrastructure/api/logouthandler"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newRestLogoutApp(factory *testrepositoryfactory.Factory, token string) *fiber.App {
	application := fiber.New()
	application.Use(stashToken(token))
	logout := logouthandler.New(factory.GetSessionRepository())
	application.Delete("/api/v1/auth/logout", logout.Logout)
	return application
}

func stashToken(token string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if token != "" {
			ctx.Locals(handler.SessionTokenKey, token)
		}
		return ctx.Next()
	}
}

func TestRestLogoutShould(t *testing.T) {
	t.Run("propagate error when session deletion fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "test-session-token"
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, token).Return(errors.New("db error"))
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		response, err := newRestLogoutApp(factory, token).Test(req, -1)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
	t.Run("succeed when no session token is present", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, "").Return(nil)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		response, err := newRestLogoutApp(factory, "").Test(req, -1)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("terminate the stashed session and return session closed message", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "rest-session-token"
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, token).Return(nil)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		response, err := newRestLogoutApp(factory, token).Test(req, -1)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		var responseData map[string]string
		data := make([]byte, response.ContentLength)
		_, _ = response.Body.Read(data)
		_ = json.Unmarshal(data, &responseData)
		assert.Equal(t, "session closed", responseData["message"])
	})
}
