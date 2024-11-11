package category_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestRestSuccess(t *testing.T) {
	factory := testrepositoryfactory.New(t)
	repository := factory.GetCategoryRepositoryMock()
	repository.On("Add", mock.Anything, mock.Anything).Return(nil)
	app := api.NewFiberApplication(factory)
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
}
