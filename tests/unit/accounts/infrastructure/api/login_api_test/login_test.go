package login_api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestRest(t *testing.T) {

	t.Run("non existing email", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := app.NewFiberApplication(factory, factory)
		user := userbuilder.New().Build()
		email := "some@email.com"
		body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, user.Password())
		var responseData map[string]string
		repository := factory.GetUserRepositoryMock()
		repository.EXPECT().GetByEmail(mock.Anything, email).Return(nil, nil)
		requestBuilder := builders.NewRequestBuilder()
		requestBuilder.WithPath("/api/v1/auth/login").WithPayload(body).WithResponseData(&responseData)
		response := getResponseFromRequestBuilder(application, requestBuilder)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "email or password are wrong", responseData["error"])
		assert.Empty(t, responseData["token"])
		repository.AssertCalled(t, "GetByEmail", mock.Anything, email)
	})

	t.Run("login with wrong data", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := app.NewFiberApplication(factory, factory)
		user := userbuilder.New().Build()
		testCases := []struct {
			name          string
			body          string
			expectedError string
		}{
			{
				"when email is missing",
				fmt.Sprintf(`{"password": "%s"}`, user.Password()),
				"email is required",
			},
			{
				"when email is empty",
				fmt.Sprintf(`{"email": "", "password": "%s"}`, user.Password()),
				"email is required",
			},
			{
				"when password is missing",
				fmt.Sprintf(`{"email": "%s"}`, user.Email()),
				"password is required",
			},
			{
				"when password is empty",
				fmt.Sprintf(`{"email": "%s", "password": ""}`, user.Email()),
				"password is required",
			},
			{
				"when body is wrong",
				fmt.Sprintf(`{"email": "%s" "password": ""}`, user.Email()),
				"wrong body",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				var responseData map[string]string
				repository := factory.GetUserRepositoryMock()
				requestBuilder := builders.NewRequestBuilder()
				requestBuilder.
					WithPath("/api/v1/auth/login").
					WithPayload(testCase.body).
					WithResponseData(&responseData)
				response := getResponseFromRequestBuilder(application, requestBuilder)

				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Equal(t, testCase.expectedError, responseData["error"])
				assert.Empty(t, responseData["token"])
				repository.AssertNotCalled(t, "GetByEmail")
			})
		}
	})

	t.Run("when email and password are correct", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		application := app.NewFiberApplication(factory, factory)
		user := userbuilder.New().Build()
		body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email(), user.Password())
		var responseData map[string]string
		userRepositoryMock := factory.GetUserRepositoryMock()
		userRepositoryMock.EXPECT().GetByEmail(mock.Anything, user.Email()).Return(user, nil)
		sessionRepositoryMock := factory.GetSessionRepositoryMock()
		sessionRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		requestBuilder := builders.NewRequestBuilder()
		requestBuilder.WithPath("/api/v1/auth/login").WithPayload(body).WithResponseData(&responseData)
		response := getResponseFromRequestBuilder(application, requestBuilder)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Empty(t, responseData["error"])
		assert.NotEmpty(t, responseData["token"])
		userRepositoryMock.AssertCalled(t, "GetByEmail", mock.Anything, user.Email())
	})
}

func getResponseFromRequestBuilder(application app.Application, requestBuilder *builders.RequestBuilder) *http.Response {
	response, err := application.Test(requestBuilder.Build())
	if err != nil {
		panic(fmt.Errorf("error making request: %w", err))
	}
	defer response.Body.Close()
	if requestBuilder.ResponseData() != nil {
		data := make([]byte, response.ContentLength)
		response.Body.Read(data)
		err = json.Unmarshal(data, requestBuilder.ResponseData())
		if err != nil {
			panic(fmt.Errorf("error unmarshalling response body: %w", err))
		}
	}
	return response
}
