package loginhandler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type loginHandler struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	authHasher        authhasher.AuthHasher
}

func New(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, authHasher authhasher.AuthHasher) loginHandler {
	return loginHandler{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		authHasher:        authHasher,
	}
}

func (self loginHandler) Login(ctx *fiber.Ctx) error {
	var body LoginBody
	if err := ctx.BodyParser(&body); err != nil {
		return odinerrors.NewErrorBuilder("wrong body").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if err := body.Validate(); err != nil {
		return err
	}
	return self.login(ctx, &body)
}

func (self loginHandler) login(ctx *fiber.Ctx, body *LoginBody) error {
	starter := sessionstarter.New(
		strings.Clone(body.Email),
		strings.Clone(body.AuthHash),
		self.userRepository,
		self.sessionRepository,
		self.authHasher,
	)
	session, user, err := starter.Start(ctx.Context())
	if err != nil {
		return err
	}
	ctx.Status(http.StatusCreated)
	return ctx.JSON(loginResponse{
		Token:              session.Token(),
		EncryptedMasterKey: user.EncryptedMasterKey(),
		KeyParams:          newKeyParamsResponse(user.KeyParams()),
	})
}

type loginResponse struct {
	Token              string             `json:"token,omitempty"`
	EncryptedMasterKey string             `json:"encrypted_master_key,omitempty"`
	KeyParams          *keyParamsResponse `json:"key_params,omitempty"`
	Error              string             `json:"error,omitempty"`
}

type keyParamsResponse struct {
	Algorithm   string `json:"algorithm"`
	Iterations  int    `json:"iterations"`
	Memory      int    `json:"memory"`
	Parallelism int    `json:"parallelism"`
	Salt        string `json:"salt"`
}

func newKeyParamsResponse(params keyparams.KeyParams) *keyParamsResponse {
	return &keyParamsResponse{
		Algorithm:   params.Algorithm(),
		Iterations:  params.Iterations(),
		Memory:      params.Memory(),
		Parallelism: params.Parallelism(),
		Salt:        params.Salt(),
	}
}
