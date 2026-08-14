package handler

import (
	"errors"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func ExternalOrFallback(err error, fallback string) string {
	var odinError *odinerrors.Error
	if errors.As(err, &odinError) && odinError.ExternalError() != "" {
		return odinError.ExternalError()
	}
	return fallback
}
