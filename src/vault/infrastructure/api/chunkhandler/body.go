package chunkhandler

import "raiseexception.dev/odin/src/shared/domain/odinerrors"

type CreateChunkBody struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (self CreateChunkBody) Validate() error {
	if self.ID == "" {
		return odinerrors.NewErrorBuilder("id is required").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if self.Content == "" {
		return odinerrors.NewErrorBuilder("content is required").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	return nil
}
