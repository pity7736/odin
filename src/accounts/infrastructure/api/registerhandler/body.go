package registerhandler

import "raiseexception.dev/odin/src/shared/domain/odinerrors"

type RegisterBody struct {
	Email              string        `json:"email"`
	AuthHash           string        `json:"auth_hash"`
	EncryptedMasterKey string        `json:"encrypted_master_key"`
	KeyParams          keyParamsBody `json:"key_params"`
}

func (self RegisterBody) Validate() error {
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
	if self.EncryptedMasterKey == "" {
		return odinerrors.NewErrorBuilder("encrypted master key is required").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	return nil
}

type keyParamsBody struct {
	Algorithm   string `json:"algorithm"`
	Salt        string `json:"salt"`
	Iterations  int    `json:"iterations"`
	Memory      int    `json:"memory"`
	Parallelism int    `json:"parallelism"`
}
