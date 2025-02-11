package sessionmodel

import "raiseexception.dev/odin/src/shared/utils"

const tokenLength uint8 = 50

type Session struct {
	token  string
	userID string
}

func New(userID string) *Session {
	return &Session{token: utils.RandomString(tokenLength), userID: userID}
}

func (self *Session) Token() string {
	return self.token
}

func (self *Session) UserID() string {
	return self.userID
}
