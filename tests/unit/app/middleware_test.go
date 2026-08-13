package app_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApp(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory, factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func TestCookieMiddlewareShould(t *testing.T) {
	t.Run("redirect to login for unauthenticated HTML request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("GET", "/accounts", nil)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.Contains(t, response.Header.Get("Location"), "/auth/login")
	})
	t.Run("return 500 when session repository returns error for cookie", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, "bad-token").Return(nil, errors.New("db error"))
		req := httptest.NewRequest("GET", "/accounts", nil)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: "bad-token"})
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
	t.Run("redirect to login when session not found", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		notFoundErr := odinerrors.NewErrorBuilder("session not found").WithTag(odinerrors.NotFound).Build()
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, "unknown-token").Return(nil, notFoundErr)
		req := httptest.NewRequest("GET", "/accounts", nil)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: "unknown-token"})
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.Contains(t, response.Header.Get("Location"), "/auth/login")
	})
	t.Run("set authenticated context and extend session for valid cookie", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		session := sessionmodel.NewFromRepository(token, "user-id", time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, token).Return(session, nil)
		factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, session).Return(nil)
		factory.GetAccountRepositoryMock().EXPECT().GetAll(mock.Anything).Return([]*accountmodel.Account{}, nil)
		req := httptest.NewRequest("GET", "/accounts", nil)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: token})
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

func TestLogoutShould(t *testing.T) {
	t.Run("HTMX logout clears cookie and redirects to login", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		session := sessionmodel.NewFromRepository(token, "user-id", time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, token).Return(session, nil)
		factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, session).Return(nil)
		factory.GetSessionRepositoryMock().EXPECT().Delete(mock.Anything, token).Return(nil)
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: token})
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "/auth/login", response.Header.Get("HX-Redirect"))
	})
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
}

func TestBearerMiddlewareShould(t *testing.T) {
	t.Run("return 401 for unauthenticated API request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("return 500 when session repository returns error for Bearer token", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, "bad-token").Return(nil, errors.New("db error"))
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
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
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
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
		factory.GetCategoryRepositoryMock().EXPECT().GetAll(mock.Anything, "user-id").Return([]*categorymodel.Category{})
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApp(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}
