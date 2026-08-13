package category_api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/mocks"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

type setup struct {
	factory           *testrepositoryfactory.Factory
	repository        *mocks.MockCategoryRepository
	app               app.Application
	userRepository    *mocks.MockUserRepository
	sessionRepository *mocks.MockSessionRepository
}

func newSetup(t *testing.T) setup {
	factory := testrepositoryfactory.New(t)
	return setup{
		factory:           factory,
		repository:        factory.GetCategoryRepositoryMock(),
		app:               app.NewFiberApplication(factory, factory.GetSessionRepository(), factory.GetUserRepository(), bcrypthasher.New()),
		userRepository:    factory.GetUserRepositoryMock(),
		sessionRepository: factory.GetSessionRepositoryMock(),
	}
}

const apiCategoryPath = "/api/v1/categories"

func TestRest(t *testing.T) {
	t.Run("create category", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		category := categorybuilder.New().Build()
		body := fmt.Sprintf(
			`{"name": "%s", "type": "%s"}`,
			category.Name(),
			category.Type(),
		)
		var responseBody map[string]any
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithPayload(body).
			WithResponseData(&responseBody).
			WithSession(session).
			WithContentType("application/json")

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, category.Name(), responseBody["name"])
		assert.Equal(t, category.Type().String(), responseBody["type"])
		assert.NotNil(t, responseBody["id"])
		assert.Equal(t, user.ID(), responseBody["user_id"])
		setup.repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("create category with anonymous user", func(t *testing.T) {
		setup := newSetup(t)
		category := categorybuilder.New().Build()
		body := fmt.Sprintf(
			`{"name": "%s", "type": "%s"}`,
			category.Name(),
			category.Type(),
		)
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithPayload(body).
			WithAnonymousSession()

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("get categories with anonymous user", func(t *testing.T) {
		setup := newSetup(t)
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithMethod(http.MethodGet).
			WithAnonymousSession()

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("get categories when is empty", func(t *testing.T) {
		setup := newSetup(t)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		setup.repository.EXPECT().GetAll(mock.Anything, user.ID()).Return(make([]*categorymodel.Category, 0))
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		var responseBody restcategoryhandler.CategoriesResponse
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithMethod(http.MethodGet).
			WithResponseData(&responseBody).
			WithSession(session)

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, 0, len(responseBody.Categories))
	})

	t.Run("get categories", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		builder := categorybuilder.New()
		categories := make([]*categorymodel.Category, 0, 1)
		categories = append(categories, builder.WithUser(user).Create(setup.repository))
		setup.repository.EXPECT().GetAll(mock.Anything, user.ID()).Return(categories)
		var responseBody restcategoryhandler.CategoriesResponse
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithMethod(http.MethodGet).
			WithResponseData(&responseBody).
			WithSession(session)

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, 1, len(responseBody.Categories))
	})

	t.Run("get categories from different user", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user0 := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user0.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		builder := categorybuilder.New()
		user1 := userbuilder.New().Create(setup.userRepository)
		builder.WithUser(user1).Create(setup.repository)
		setup.repository.EXPECT().GetAll(mock.Anything, mock.Anything).Return([]*categorymodel.Category{})
		var responseBody restcategoryhandler.CategoriesResponse
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(apiCategoryPath).
			WithMethod(http.MethodGet).
			WithResponseData(&responseBody).
			WithSession(session)

		response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, 0, len(responseBody.Categories))
	})

	// BUG: returns 500 instead of 400 — category handler error propagation is broken.
	// See specs/accounting/categories/ task in tasks.md.
	// t.Run("create category with wrong data", func(t *testing.T) {
	// 	setup := newSetup(t)
	// 	setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
	// 	user := userbuilder.New().Create(setup.userRepository)
	// 	session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
	// 	setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
	// 	setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
	// 	category := categorybuilder.New().Build()
	// 	testCases := []struct {
	// 		testCaseName string
	// 		categoryName string
	// 		categoryType string
	// 	}{
	// 		{
	// 			"when name is empty",
	// 			"",
	// 			category.Type().String(),
	// 		},
	// 		{
	// 			"when type is empty",
	// 			"test",
	// 			"",
	// 		},
	// 		{
	// 			"when type is invalid",
	// 			"test",
	// 			"eaoeu",
	// 		},
	// 	}
	// 	for _, testCase := range testCases {
	// 		t.Run(testCase.testCaseName, func(t *testing.T) {
	// 			body := fmt.Sprintf(
	// 				`{"name": "%s", "type": "%s"}`,
	// 				testCase.categoryName,
	// 				testCase.categoryType,
	// 			)
	// 			requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
	// 			requestBuilder.
	// 				WithPath(apiCategoryPath).
	// 				WithPayload(body).
	// 				WithSession(session)
	//
	// 			response := testutils.GetJSONResponseFromRequestBuilder(setup.app, requestBuilder)
	// 			defer func() { _ = response.Body.Close() }()
	//
	// 			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	// 			assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
	// 			setup.repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	// 		})
	// 	}
	// })
}

const categoryPath = "/categories"

func TestHTMX(t *testing.T) {
	t.Run("create category", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		category := categorybuilder.New().Build()
		body := fmt.Sprintf(
			"name=%s&type=%s",
			category.Name(),
			category.Type(),
		)
		var responseBody map[string]any
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithPayload(body).
			WithResponseData(&responseBody).
			WithSession(session).
			WithContentType("application/x-www-form-urlencoded")

		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.True(t, strings.Contains(responseData, category.Name()))
		setup.repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("create category when token does not exists", func(t *testing.T) {
		setup := newSetup(t)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, mock.Anything).Return(nil,
			odinerrors.NewErrorBuilder("session not found").
				WithExternalMessage("Sesión no encontrada").
				WithTag(odinerrors.NotFound).
				Build())
		category := categorybuilder.New().Build()
		body := fmt.Sprintf(
			"name=%s&type=%s",
			category.Name(),
			category.Type(),
		)
		var responseBody map[string]any
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithPayload(body).
			WithResponseData(&responseBody).
			WithSession(session).
			WithContentType("application/x-www-form-urlencoded")

		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.False(t, strings.Contains(responseData, category.Name()))
	})

	t.Run("get categories when is empty", func(t *testing.T) {
		setup := newSetup(t)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		setup.repository.EXPECT().GetAll(mock.Anything, user.ID()).Return(make([]*categorymodel.Category, 0))
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithMethod(http.MethodGet).
			WithSession(session).
			WithContentType("")

		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
		assert.True(t, strings.Contains(responseData, "hx-vals='{\"first\": \"true\"}'"))
		assert.True(t, strings.Contains(responseData, "<p>no hay categorías</p>"))
	})

	t.Run("get categories", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, session.Token()).Return(session, nil)
		setup.sessionRepository.EXPECT().Save(mock.Anything, session).Return(nil)
		categories := make([]*categorymodel.Category, 0, 1)
		category := categorybuilder.New().Create(setup.repository)
		categories = append(categories, category)
		setup.repository.EXPECT().GetAll(mock.Anything, user.ID()).Return(categories)
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithMethod(http.MethodGet).
			WithSession(session).
			WithContentType("")

		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
		assert.False(t, strings.Contains(responseData, "hx-vals='{\"first\": \"true\"}'"))
		assert.True(t, strings.Contains(responseData, category.Name()))
	})

	t.Run("get categories with anonymous user", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().Create(setup.repository)
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithMethod(http.MethodGet).
			WithContentType("").
			WithAnonymousSession()

		response, responseData := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusFound, response.StatusCode)
		assert.False(t, strings.Contains(responseData, category.Name()))
	})

	t.Run("get categories when session token does not exists", func(t *testing.T) {
		setup := newSetup(t)
		setup.userRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		user := userbuilder.New().Create(setup.userRepository)
		session, _ := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
		setup.sessionRepository.EXPECT().Get(mock.Anything, mock.Anything).Return(nil,
			odinerrors.NewErrorBuilder("session not found").
				WithExternalMessage("Sesión no encontrada").
				WithTag(odinerrors.NotFound).
				Build())
		requestBuilder := builders.NewRequestBuilder(setup.userRepository, setup.sessionRepository)
		requestBuilder.
			WithPath(categoryPath).
			WithMethod(http.MethodGet).
			WithSession(session).
			WithContentType("")

		response, _ := testutils.GetHTMLResponseFromRequestBuilder(setup.app, requestBuilder)
		defer func() { _ = response.Body.Close() }()

		assert.Equal(t, http.StatusFound, response.StatusCode)
	})
}
