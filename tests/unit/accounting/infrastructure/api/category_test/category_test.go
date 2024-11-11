package category_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestRestSuccess(t *testing.T) {
	factory := testrepositoryfactory.New(t)
	repository := factory.GetCategoryRepositoryMock()
	repository.On("Add", mock.Anything, mock.Anything).Return(nil)
	app := api.NewFiberApplication(factory)
	body := `{"name": "test", "type": "EXPENSE", "user": "test@raiseexception.dev"}`
	request := httptest.NewRequest("POST", "/v1/categories", bytes.NewReader([]byte(body)))
	request.Header.Add("Content-Type", "application/json")

	response, _ := app.Test(request)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	repository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
}
