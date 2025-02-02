package sessionstarter

import (
	"context"
	"errors"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
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

func (self *SessionStarter) Start(ctx context.Context) (*sessionmodel.Session, error) {
	user, err := self.userRepository.GetByEmail(ctx, self.email)
	if err != nil {
		return nil, err
	}
	return self.start(ctx, user)
}

func (self *SessionStarter) start(ctx context.Context, user *usermodel.User) (*sessionmodel.Session, error) {
	if user != nil && user.CheckPassword(self.password) {
		return self.createSession(ctx)
	}
	return nil, errors.New("email or password are wrong")
}

func (self *SessionStarter) createSession(ctx context.Context) (*sessionmodel.Session, error) {
	err := self.sessionRepository.Add(ctx, sessionmodel.New())
	if err != nil {
		return nil, err
	}
	return sessionmodel.New(), nil
}
