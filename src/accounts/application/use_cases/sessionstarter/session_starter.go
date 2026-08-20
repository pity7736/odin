package sessionstarter

import (
	"context"
	"errors"

	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type SessionStarter struct {
	email             string
	authHash          string
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	authHasher        authhasher.AuthHasher
}

func New(
	email, authHash string,
	userRepository repositories.UserRepository,
	sessionRepository repositories.SessionRepository,
	authHasher authhasher.AuthHasher,
) SessionStarter {

	return SessionStarter{
		email:             email,
		authHash:          authHash,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		authHasher:        authHasher,
	}
}

func (self SessionStarter) Start(ctx context.Context) (*sessionmodel.Session, *usermodel.User, error) {
	user, err := self.userRepository.GetByEmail(ctx, self.email)
	if err != nil {
		var odinError *odinerrors.Error
		if errors.As(err, &odinError) && odinError.Tag() == odinerrors.NotFound {
			return nil, nil, wrongCredentials()
		}
		return nil, nil, err
	}
	return self.start(ctx, user)
}

func (self SessionStarter) start(ctx context.Context, user *usermodel.User) (*sessionmodel.Session, *usermodel.User, error) {
	if self.authHasher.Compare(user.AuthHashDigest(), self.authHash) {
		return self.createSession(ctx, user)
	}
	return nil, nil, wrongCredentials()
}

func (self SessionStarter) createSession(ctx context.Context, user *usermodel.User) (*sessionmodel.Session, *usermodel.User, error) {
	session, err := sessionmodel.New(user.ID(), sessionmodel.DefaultTTL)
	if err != nil {
		return nil, nil, err
	}
	err = self.sessionRepository.Add(ctx, session)
	if err != nil {
		return nil, nil, err
	}
	return session, user, nil
}

func wrongCredentials() error {
	return odinerrors.NewErrorBuilder("email or password are wrong").
		WithExternalMessage("Correo o contraseña incorrectos").
		WithTag(odinerrors.Unauthorized).
		Build()
}
