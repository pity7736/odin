package sessionmodel

import "raiseexception.dev/odin/src/shared/utils"

const tokenLength uint8 = 50

type Session struct {
	token string
}

func New() *Session {
	return &Session{token: utils.RandomString(tokenLength)}
}

func (s *Session) Token() string {
	return s.token
}
