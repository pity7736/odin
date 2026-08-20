package chunkmodel

import (
	"github.com/google/uuid"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type EncryptedChunk struct {
	id      string
	ownerID string
	content string
}

func New(id, ownerID, content string) (*EncryptedChunk, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, odinerrors.NewErrorBuilder("id must be a valid uuid").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if ownerID == "" {
		return nil, odinerrors.NewErrorBuilder("owner id is required").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if content == "" {
		return nil, odinerrors.NewErrorBuilder("content is required").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	return &EncryptedChunk{id: id, ownerID: ownerID, content: content}, nil
}

func (self *EncryptedChunk) ID() string {
	return self.id
}

func (self *EncryptedChunk) OwnerID() string {
	return self.ownerID
}

func (self *EncryptedChunk) Content() string {
	return self.content
}
