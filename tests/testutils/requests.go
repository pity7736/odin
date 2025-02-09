package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
)

func GetHtmlResponseFromRequestBuilder(application app.Application, requestBuilder *builders.RequestBuilder) (*http.Response, string) {
	response, err := application.Test(requestBuilder.Build())
	if err != nil {
		panic(fmt.Errorf("error making request: %w", err))
	}
	defer response.Body.Close()
	data := make([]byte, response.ContentLength)
	response.Body.Read(data)
	return response, string(data)
}

func GetJsonResponseFromRequestBuilder(application app.Application, requestBuilder *builders.RequestBuilder) *http.Response {
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
