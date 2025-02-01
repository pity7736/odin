package sessionstarter

import (
	"context"
	"errors"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/accountsrepositoryfactory"
)

type SessionStarter struct {
	email             string
	password          string
	sessionRepository repositories.SessionRepository
	userRepository    repositories.UserRepository
}

func New(email, password string,
	factory accountsrepositoryfactory.AccountsRepositoryFactory) *SessionStarter {
	return &SessionStarter{
		email:             email,
		password:          password,
		sessionRepository: factory.GetSessionRepository(),
		userRepository:    factory.GetUserRepository(),
	}
}

func (self *SessionStarter) Start(ctx context.Context) (string, error) {
	user, err := self.userRepository.GetByEmail(ctx, self.email)
	if err != nil {
		return "", err
	}
	return self.start(ctx, user)
}

func (self *SessionStarter) start(ctx context.Context, user *usermodel.User) (string, error) {
	if user != nil && user.CheckPassword(self.password) {
		return self.createToken(ctx)
	}
	return "", errors.New("email or password are wrong")
}

func (self *SessionStarter) createToken(ctx context.Context) (string, error) {
	token := "token" // TODO: create a real token
	err := self.sessionRepository.Add(ctx, token)
	if err != nil {
		return "", err
	}
	return token, nil
}
