package chunkhandler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	"raiseexception.dev/odin/src/vault/application/use_cases/chunkcreator"
	"raiseexception.dev/odin/src/vault/domain/repositories"
)

type chunkHandler struct {
	chunkRepository repositories.ChunkRepository
}

func New(chunkRepository repositories.ChunkRepository) chunkHandler {
	return chunkHandler{
		chunkRepository: chunkRepository,
	}
}

func (self chunkHandler) Create(ctx *fiber.Ctx) error {
	var body CreateChunkBody
	if err := ctx.BodyParser(&body); err != nil {
		return odinerrors.NewErrorBuilder("wrong body").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if err := body.Validate(); err != nil {
		return err
	}
	return self.create(ctx, &body)
}

func (self chunkHandler) create(ctx *fiber.Ctx, body *CreateChunkBody) error {
	ownerID := ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext).UserID()
	creator := chunkcreator.New(
		strings.Clone(body.ID),
		ownerID,
		strings.Clone(body.Content),
		self.chunkRepository,
	)
	chunk, err := creator.Create(ctx.Context())
	if err != nil {
		return err
	}
	ctx.Status(http.StatusCreated)
	return ctx.JSON(createChunkResponse{ID: chunk.ID()})
}

type createChunkResponse struct {
	ID string `json:"id"`
}
