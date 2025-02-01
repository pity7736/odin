package builders

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
)

type RequestBuilder struct {
	method       string
	path         string
	contentType  string
	payload      io.Reader
	responseData any
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		method:      "POST",
		path:        "/",
		contentType: "application/json",
	}
}

func (self *RequestBuilder) WithMethod(method string) *RequestBuilder {
	self.method = method
	return self
}

func (self *RequestBuilder) WithPath(path string) *RequestBuilder {
	self.path = path
	return self
}

func (self *RequestBuilder) WithPayload(payload string) *RequestBuilder {
	self.payload = bytes.NewReader([]byte(payload))
	return self
}

func (self *RequestBuilder) WithResponseData(data any) *RequestBuilder {
	self.responseData = data
	return self
}

func (self *RequestBuilder) Build() *http.Request {
	request := httptest.NewRequest(self.method, self.path, self.payload)
	request.Header.Add("Content-Type", self.contentType)
	return request
}

func (self *RequestBuilder) ResponseData() any {
	return self.responseData
}
