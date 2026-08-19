package loginhandler

import "raiseexception.dev/odin/src/shared/domain/odinerrors"

type LoginBody struct {
	Email    string `json:"email"`
	AuthHash string `json:"auth_hash"`
}

func (self LoginBody) Validate() error {
	if self.Email == "" {
		return odinerrors.NewErrorBuilder("email is required").
			WithExternalMessage("El correo es obligatorio").
			WithTag(odinerrors.Domain).
			Build()
	}
	if self.AuthHash == "" {
		return odinerrors.NewErrorBuilder("auth hash is required").
			WithExternalMessage("La contraseña es obligatoria").
			WithTag(odinerrors.Domain).
			Build()
	}
	return nil
}
