package category_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestRestSuccess(t *testing.T) {

	t.Run("create category", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		category := categorybuilder.New().WithDefaultUser().Build()
		body := fmt.Sprintf(
			`{"name": "%s", "type": "%s", "user": "%s"}`,
			category.Name(),
			category.Type(),
			category.User().Email(),
		)
		var responseBody map[string]any

		response := makeRequestAndGetResponse[map[string]any](
			setup,
			"POST",
			"/v1/categories",
			&body,
			&responseBody,
		)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, category.Name(), responseBody["name"])
		assert.Equal(t, category.Type().String(), responseBody["type"])
		assert.NotNil(t, responseBody["id"])
		assert.Equal(t, category.User().Email(), responseBody["user"])
		setup.repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("get categories when is empty", func(t *testing.T) {
		setup := newSetup(t)
		setup.repository.EXPECT().GetAll(mock.Anything).Return(make([]*category.Category, 0))
		var responseBody restcategoryhandler.CategoriesResponse
		response := makeRequestAndGetResponse[restcategoryhandler.CategoriesResponse](
			setup,
			"GET",
			"/v1/categories",
			nil,
			&responseBody,
		)

		assert.Equal(t, http.StatusOK, response.StatusCode)
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
			"/v1/categories",
			nil,
			&responseBody,
		)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 1, len(responseBody.Categories))
	})
}

func makeRequestAndGetResponse[R any](setup setup, method, path string, payload *string, responseBody *R) *http.Response {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader([]byte(*payload))
	}
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	response, _ := setup.app.Test(request)
	data := make([]byte, response.ContentLength)
	response.Body.Read(data)
	json.Unmarshal(data, &responseBody)
	return response
}
