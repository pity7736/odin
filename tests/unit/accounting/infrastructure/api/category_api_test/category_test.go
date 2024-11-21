package category_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/unit/mocks"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

type setup struct {
	factory    *testrepositoryfactory.Factory
	repository *mocks.MockCategoryRepository
	app        api.Application
}

func newSetup(t *testing.T) setup {
	factory := testrepositoryfactory.New(t)
	return setup{
		factory:    factory,
		repository: factory.GetCategoryRepositoryMock(),
		app:        api.NewFiberApplication(factory),
	}
}

const categoryPath = "/v1/categories"

func TestRest(t *testing.T) {

	t.Run("create category", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithDefaultUser().Build()
		body := fmt.Sprintf(
			`{"name": "%s", "type": "%s"}`,
			category.Name(),
			category.Type(),
		)
		var responseBody map[string]any

		response := makeRequestAndGetResponse[map[string]any](
			setup,
			"POST",
			categoryPath,
			&body,
			&responseBody,
		)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, category.Name(), responseBody["name"])
		assert.Equal(t, category.Type().String(), responseBody["type"])
		assert.NotNil(t, responseBody["id"])
		setup.repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("get categories when is empty", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().GetAll(mock.Anything).Return(make([]*category.Category, 0))
		var responseBody restcategoryhandler.CategoriesResponse
		response := makeRequestAndGetResponse[restcategoryhandler.CategoriesResponse](
			setup,
			"GET",
			categoryPath,
			nil,
			&responseBody,
		)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, 0, len(responseBody.Categories))
	})

	t.Run("get categories", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		builder := categorybuilder.New()
		categories := make([]*category.Category, 0, 1)
		categories = append(categories, builder.WithDefaultUser().Create(setup.repository))
		setup.repository.EXPECT().GetAll(mock.Anything).Return(categories)
		var responseBody restcategoryhandler.CategoriesResponse
		response := makeRequestAndGetResponse[restcategoryhandler.CategoriesResponse](
			setup,
			"GET",
			categoryPath,
			nil,
			&responseBody,
		)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
		assert.Equal(t, 1, len(responseBody.Categories))
	})

	t.Run("create category with wrong data", func(t *testing.T) {
		setup := newSetup(t)
		category := categorybuilder.New().WithDefaultUser().Build()
		testCases := []struct {
			testCaseName string
			categoryName string
			categoryType string
		}{
			{
				"when name is empty",
				"",
				category.Type().String(),
			},
			{
				"when type is empty",
				"test",
				"",
			},
			{
				"when type is invalid",
				"test",
				"eaoeu",
			},
		}
		for _, testCase := range testCases {
			// TODO: send appropriate error message
			t.Run(testCase.testCaseName, func(t *testing.T) {
				body := fmt.Sprintf(
					`{"name": "%s", "type": "%s"}`,
					testCase.categoryName,
					testCase.categoryType,
				)
				var responseBody map[string]any

				response := makeRequestAndGetResponse[map[string]any](
					setup,
					"POST",
					categoryPath,
					&body,
					&responseBody,
				)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Equal(t, fiber.MIMEApplicationJSON, response.Header.Get("content-type"))
				setup.repository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
			})
		}
	})
}

func makeRequestAndGetResponse[R any](setup setup, method, path string, payload *string, responseBody *R) *http.Response {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader([]byte(*payload))
	}
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	userID, _ := uuid.NewV7()
	request.Header.Add("Authorization", userID.String())
	response, _ := setup.app.Test(request)
	data := make([]byte, response.ContentLength)
	response.Body.Read(data)
	json.Unmarshal(data, &responseBody)
	return response
}

func TestHTMX(t *testing.T) {

	t.Run("get categories when is empty", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().GetAll(mock.Anything).Return(make([]*category.Category, 0))
		request := httptest.NewRequest("GET", categoryPath, nil)

		response, _ := setup.app.Test(request)
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		responseBody := string(data)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
		assert.True(t, strings.Contains(responseBody, "hx-vals='{\"first\": \"true\"}'"))
		assert.True(t, strings.Contains(responseBody, "<p>no hay categorías</p>"))
	})

	t.Run("get categories", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		builder := categorybuilder.New()
		categories := make([]*category.Category, 0, 1)
		category := builder.WithDefaultUser().Create(setup.repository)
		categories = append(categories, category)
		setup.repository.EXPECT().GetAll(mock.Anything).Return(categories)
		request := httptest.NewRequest("GET", categoryPath, nil)

		response, _ := setup.app.Test(request)
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		responseBody := string(data)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
		assert.False(t, strings.Contains(responseBody, "hx-vals='{\"first\": \"true\"}'"))
		assert.True(t, strings.Contains(responseBody, category.Name()))
	})
}
