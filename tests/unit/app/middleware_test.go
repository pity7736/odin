package app_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApp(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func TestLogoutShould(t *testing.T) {
	t.Run("REST logout returns session closed message", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		session := sessionmodel.NewFromRepository(token, "user-id", time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, token).Return(session, nil)
		factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, session).Return(nil)
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, token).Return(nil)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("REST logout terminates the bearer session even when a session cookie is also present", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		bearerToken := "bearer-token"
		bearerSession := sessionmodel.NewFromRepository(bearerToken, "user-id", time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, bearerToken).Return(bearerSession, nil)
		factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, bearerSession).Return(nil)
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, bearerToken).Return(nil)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+bearerToken)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: "cookie-token"})
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

func TestBearerMiddlewareShould(t *testing.T) {
	t.Run("return 401 for unauthenticated API request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("return 500 when session repository returns error for Bearer token", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, "bad-token").Return(nil, errors.New("db error"))
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
	t.Run("treat as anonymous when Bearer token is not found", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		notFoundErr := odinerrors.NewErrorBuilder("session not found").WithTag(odinerrors.NotFound).Build()
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, "unknown-token").Return(nil, notFoundErr)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer unknown-token")
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("set authenticated context and extend session for valid Bearer token", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-bearer-token"
		session := sessionmodel.NewFromRepository(token, "user-id", time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, token).Return(session, nil)
		factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, session).Return(nil)
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, token).Return(nil)
		req := httptest.NewRequest("DELETE", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}
