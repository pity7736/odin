package sessionstarter

import (
	"context"

	"raiseexception.dev/odin/src/accounts/application/passwordhasher"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type SessionStarter struct {
	email             string
	password          string
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	passwordHasher    passwordhasher.PasswordHasher
}

func New(
	email, password string,
	userRepository repositories.UserRepository,
	sessionRepository repositories.SessionRepository,
	passwordHasher passwordhasher.PasswordHasher,
) *SessionStarter {

	return &SessionStarter{
		email:             email,
		password:          password,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		passwordHasher:    passwordHasher,
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
	if user != nil && self.passwordHasher.Compare(user.HashedPassword(), self.password) {
		return self.createSession(ctx, user)
	}
	return nil, odinerrors.NewErrorBuilder("email or password are wrong").
		WithExternalMessage("Correo o contraseña incorrectos").
		WithTag(odinerrors.DOMAIN).
		Build()
}

func (self *SessionStarter) createSession(ctx context.Context, user *usermodel.User) (*sessionmodel.Session, error) {
	session, err := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
	if err != nil {
		return nil, err
	}
	err = self.sessionRepository.Add(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
