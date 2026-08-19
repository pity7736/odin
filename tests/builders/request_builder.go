package builders

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/tests/builders/userbuilder"
)

type RequestBuilder struct {
	method            string
	path              string
	contentType       string
	payload           io.Reader
	responseData      any
	session           *sessionmodel.Session
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	withSession       bool
	user              *usermodel.User
}

func NewRequestBuilder(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository) *RequestBuilder {
	return &RequestBuilder{
		method:            "POST",
		path:              "/",
		contentType:       "",
		withSession:       true,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
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
	self.withSession = true
	return self
}

func (self *RequestBuilder) WithAnonymousSession() *RequestBuilder {
	self.withSession = false
	return self
}

func (self *RequestBuilder) WithUser(user *usermodel.User) *RequestBuilder {
	self.user = user
	return self
}

func (self *RequestBuilder) Build() *http.Request {
	request := httptest.NewRequest(self.method, self.path, self.payload)
	if self.contentType != "" {
		request.Header.Set("Content-Type", self.contentType)
	}
	if self.withSession {
		session := self.session
		if session == nil {
			user := self.user
			if user == nil {
				user = userbuilder.New().Create(self.userRepository)
			}
			sessionStarter := sessionstarter.New(
				user.Email(),
				userbuilder.DefaultPassword,
				self.userRepository,
				self.sessionRepository,
				bcrypthasher.New(),
			)
			session, _, _ = sessionStarter.Start(context.TODO())
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Token()))
	}
	return request
}

func (self *RequestBuilder) ResponseData() any {
	return self.responseData
}

func (self *RequestBuilder) IsContentType(contentType string) bool {
	return contentType == self.contentType
}
