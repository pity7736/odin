package builders

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
)

type RequestBuilder struct {
	method       string
	path         string
	contentType  string
	payload      io.Reader
	responseData any
	session      *sessionmodel.Session
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		method:      "POST",
		path:        "/",
		contentType: "",
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

func (self *RequestBuilder) WithContentType(contentType string) *RequestBuilder {
	self.contentType = contentType
	return self
}

func (self *RequestBuilder) WithSession(session *sessionmodel.Session) *RequestBuilder {
	self.session = session
	return self
}

func (self *RequestBuilder) Build() *http.Request {
	request := httptest.NewRequest(self.method, self.path, self.payload)
	if self.contentType != "" {
		request.Header.Set("Content-Type", self.contentType)
	}
	if self.session != nil {
		if self.contentType == "application/json" {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.session.Token()))
		} else {
			request.AddCookie(&http.Cookie{
				Name:     "__Secure-odin-session",
				Value:    self.session.Token(),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
		}
	}
	return request
}

func (self *RequestBuilder) ResponseData() any {
	return self.responseData
}

func (self *RequestBuilder) IsContentType(contentType string) bool {
	return contentType == self.contentType
}
