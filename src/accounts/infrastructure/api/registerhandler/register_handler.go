package registerhandler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/application/use_cases/userregistrar"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type registerHandler struct {
	userRepository repositories.UserRepository
	authHasher     authhasher.AuthHasher
}

func New(userRepository repositories.UserRepository, authHasher authhasher.AuthHasher) registerHandler {
	return registerHandler{
		userRepository: userRepository,
		authHasher:     authHasher,
	}
}

func (self registerHandler) Register(ctx *fiber.Ctx) error {
	var body RegisterBody
	if err := ctx.BodyParser(&body); err != nil {
		return odinerrors.NewErrorBuilder("wrong body").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if err := body.Validate(); err != nil {
		return err
	}
	return self.register(ctx, &body)
}

func (self registerHandler) register(ctx *fiber.Ctx, body *RegisterBody) error {
	keyParams, err := keyparams.New(
		strings.Clone(body.KeyParams.Algorithm),
		body.KeyParams.Iterations,
		body.KeyParams.Memory,
		body.KeyParams.Parallelism,
		strings.Clone(body.KeyParams.Salt),
	)
	if err != nil {
		return err
	}
	registrar := userregistrar.New(
		strings.Clone(body.Email),
		strings.Clone(body.AuthHash),
		strings.Clone(body.EncryptedMasterKey),
		keyParams,
		self.userRepository,
		self.authHasher,
	)
	user, err := registrar.Register(ctx.Context())
	if err != nil {
		return err
	}
	ctx.Status(http.StatusCreated)
	return ctx.JSON(registerResponse{
		ID:    user.ID(),
		Email: user.Email(),
	})
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
