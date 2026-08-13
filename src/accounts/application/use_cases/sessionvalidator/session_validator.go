package sessionvalidator

import (
	"context"
	"errors"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type SessionValidator struct {
	sessionRepository repositories.SessionRepository
}

func New(sessionRepository repositories.SessionRepository) SessionValidator {
	return SessionValidator{sessionRepository: sessionRepository}
}

func (self SessionValidator) Validate(ctx context.Context, token string) (*sessionmodel.Session, error) {
	session, err := self.sessionRepository.Get(ctx, token)
	if err != nil {
		var odinError *odinerrors.Error
		if errors.As(err, &odinError) && (odinError.Tag() == odinerrors.NotFound || odinError.Tag() == odinerrors.Domain) {
			return nil, nil
		}
		return nil, err
	}
	session.Extend(sessionmodel.DefaultTTL)
	if err := self.sessionRepository.Save(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}
