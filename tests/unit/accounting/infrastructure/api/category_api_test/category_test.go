package category_api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/rest/restcategoryhandler"
	"raiseexception.dev/odin/src/shared/domain/user"
	"raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestRestSuccess(t *testing.T) {

	t.Run("create category", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetCategoryRepositoryMock()
		app := api.NewFiberApplication(factory)
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		categoryName := "test"
		categoryType := constants.EXPENSE.String()
		user := "test@raiseexception.dev"
		body := fmt.Sprintf(
			`{"name": "%s", "type": "%s", "user": "%s"}`,
			categoryName,
			categoryType,
			user,
		)
		request := httptest.NewRequest("POST", "/v1/categories", bytes.NewReader([]byte(body)))
		request.Header.Add("Content-Type", "application/json")

		response, _ := app.Test(request)

		var responseBody map[string]any
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		json.Unmarshal(data, &responseBody)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, categoryName, responseBody["name"])
		assert.Equal(t, categoryType, responseBody["type"])
		assert.NotNil(t, responseBody["id"])
		assert.Equal(t, user, responseBody["user"])
		repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("get categories when is empty", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetCategoryRepositoryMock()
		app := api.NewFiberApplication(factory)
		repository.EXPECT().GetAll(mock.Anything).Return(make([]*category.Category, 0))

		request := httptest.NewRequest("GET", "/v1/categories", nil)
		request.Header.Add("Content-Type", "application/json")
		response, _ := app.Test(request)

		var responseBody restcategoryhandler.CategoriesResponse
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		json.Unmarshal(data, &responseBody)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 0, len(responseBody.Categories))
	})

	t.Run("get categories", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repository := factory.GetCategoryRepositoryMock()
		repository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		categoryName := "test"
		user := getUser()
		categoryType := constants.EXPENSE
		command := categorycommand.New(categoryName, categoryType, user)
		categoryCreator := categorycreator.New(command, repository)
		categories := make([]*category.Category, 0, 1)
		category, _ := categoryCreator.Create(context.TODO())
		categories = append(categories, category)
		app := api.NewFiberApplication(factory)
		repository.EXPECT().GetAll(mock.Anything).Return(categories)

		request := httptest.NewRequest("GET", "/v1/categories", nil)
		request.Header.Add("Content-Type", "application/json")
		response, _ := app.Test(request)

		var responseBody restcategoryhandler.CategoriesResponse
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		json.Unmarshal(data, &responseBody)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 1, len(responseBody.Categories))
	})
}

func getUser() *user.User {
	return user.New("test@raiseexception.dev")
}
