package category_api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	handler "raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newApplication(factory *testrepositoryfactory.Factory) app.Application {
	return app.NewFiberApplication(factory, factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New())
}

func setupAuthMocks(factory *testrepositoryfactory.Factory, token, userID string) *sessionmodel.Session {
	session := sessionmodel.NewFromRepository(token, userID, time.Now(), time.Now().Add(sessionmodel.DefaultTTL))
	factory.GetSessionRepositoryMock().EXPECT().Get(mock.Anything, token).Return(session, nil)
	factory.GetSessionRepositoryMock().EXPECT().Save(mock.Anything, session).Return(nil)
	return session
}

func TestRestCategoryShould(t *testing.T) {
	t.Run("return 401 for unauthenticated POST request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		body := `{"name": "Comida", "type": "expense"}`
		req := httptest.NewRequest("POST", "/api/v1/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("return 401 for unauthenticated GET request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("create category when authenticated", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "auth-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		factory.GetCategoryRepositoryMock().EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		body := `{"name": "Comida", "type": "expense"}`
		req := httptest.NewRequest("POST", "/api/v1/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("return categories when authenticated", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "auth-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		factory.GetCategoryRepositoryMock().EXPECT().GetAll(mock.Anything, userID).Return([]*categorymodel.Category{})
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("return non-empty categories list when authenticated", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "auth-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		category := categorybuilder.New().WithUserID(userID).Build()
		factory.GetCategoryRepositoryMock().EXPECT().GetAll(mock.Anything, userID).Return([]*categorymodel.Category{category})
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("return error when category type is invalid", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "auth-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		body := `{"name": "Comida", "type": "invalid_type"}`
		req := httptest.NewRequest("POST", "/api/v1/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+token)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
}

func TestHtmxCategoryShould(t *testing.T) {
	t.Run("create category when authenticated via cookie", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		factory.GetCategoryRepositoryMock().EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		body := "name=Comida&type=expense"
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationForm)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: token})
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("return categories when authenticated via cookie", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		factory.GetCategoryRepositoryMock().EXPECT().GetAll(mock.Anything, userID).Return([]*categorymodel.Category{})
		req := httptest.NewRequest("GET", "/categories", nil)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: token})
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("set HX-Refresh when creating the first category", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		token := "valid-token"
		userID := "user-id"
		setupAuthMocks(factory, token, userID)
		factory.GetCategoryRepositoryMock().EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		body := "name=Comida&type=expense&first=true"
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationForm)
		req.AddCookie(&http.Cookie{Name: handler.SessionName, Value: token})
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "true", response.Header.Get("HX-Refresh"))
	})
	t.Run("redirect to login for unauthenticated POST request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		body := "name=Comida&type=expense"
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", fiber.MIMEApplicationForm)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.Contains(t, response.Header.Get("Location"), "/auth/login")
	})
	t.Run("redirect to login for unauthenticated GET request", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		req := httptest.NewRequest("GET", "/categories", nil)
		response, err := newApplication(factory).Test(req)
		assert.Nil(t, err)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.Contains(t, response.Header.Get("Location"), "/auth/login")
	})
}
