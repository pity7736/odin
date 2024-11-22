package sessionstarter

import (
	"errors"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type SessionStarter struct {
	email             string
	password          string
	sessionRepository repositories.SessionRepository
	userRepository    repositories.UserRepository
}

func New(email, password string,
	sessionRepository repositories.SessionRepository,
	userRepository repositories.UserRepository) *SessionStarter {
	return &SessionStarter{
		email:             email,
		password:          password,
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (self *SessionStarter) Start() (string, error) {
	user, err := self.userRepository.GetByEmail(self.email)
	if err != nil {
		return "", err
	}
	return self.start(user)
}

func (self *SessionStarter) start(user *usermodel.User) (string, error) {
	if user != nil && user.CheckPassword(self.password) {
		return self.createToken()
	}
	return "", errors.New("email or password are wrong")
}

func (self *SessionStarter) createToken() (string, error) {
	token := "token"
	err := self.sessionRepository.Add(token)
	if err != nil {
		return "", err
	}
	return token, nil
}
