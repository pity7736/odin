package requestcontext

import (
	"errors"

	"github.com/google/uuid"
)

const Key = "requestContext"

type RequestContext struct {
	userID    string
	requestID string
}

func New(userID string) (*RequestContext, error) {
	if userID == "" {
		return nil, errors.New("user id cannot be empty")
	}
	requestID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &RequestContext{userID: userID, requestID: requestID.String()}, nil
}

func NewAnonymous() *RequestContext {
	requestID, _ := uuid.NewV7()
	return &RequestContext{requestID: requestID.String()}
}

func (self *RequestContext) UserID() string {
	return self.userID
}

func (self *RequestContext) RequestID() string {
	return self.requestID
}

func (self *RequestContext) IsAuthenticated() bool {
	return self.userID != ""
}
