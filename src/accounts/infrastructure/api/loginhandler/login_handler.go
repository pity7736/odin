package loginhandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type LoginResult struct {
	Token              string
	EncryptedMasterKey string
	KeyParams          keyparams.KeyParams
}

type loginHandler struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	authHasher        authhasher.AuthHasher
	ctx               *fiber.Ctx
}

func New(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, authHasher authhasher.AuthHasher, ctx *fiber.Ctx) *loginHandler {
	return &loginHandler{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		authHasher:        authHasher,
		ctx:               ctx,
	}
}

func (self *loginHandler) Login() error {
	self.ctx.Set("Content-Type", fiber.MIMEApplicationJSON)
	var body LoginBody
	if err := self.validateRequestBody(&body); err != nil {
		self.ctx.Status(http.StatusBadRequest)
		return self.renderError(err)
	}
	return self.login(&body)
}

func (self *loginHandler) validateRequestBody(body *LoginBody) error {
	if err := self.ctx.BodyParser(body); err != nil {
		return odinerrors.NewErrorBuilder("wrong body").
			WithExternalMessage("Datos de solicitud inválidos").
			WithTag(odinerrors.Domain).
			Build()
	}
	if body.Email == "" {
		return odinerrors.NewErrorBuilder("email is required").
			WithExternalMessage("El correo es obligatorio").
			WithTag(odinerrors.Domain).
			Build()
	}
	if body.AuthHash == "" {
		return odinerrors.NewErrorBuilder("auth hash is required").
			WithExternalMessage("La contraseña es obligatoria").
			WithTag(odinerrors.Domain).
			Build()
	}
	return nil
}

func (self *loginHandler) login(body *LoginBody) error {
	starter := sessionstarter.New(
		strings.Clone(body.Email),
		strings.Clone(body.AuthHash),
		self.userRepository,
		self.sessionRepository,
		self.authHasher,
	)
	session, user, err := starter.Start(self.ctx.Context())
	if err != nil {
		return self.handleError(err)
	}
	self.ctx.Status(http.StatusCreated)
	result := &LoginResult{
		Token:              session.Token(),
		EncryptedMasterKey: user.EncryptedMasterKey(),
		KeyParams:          user.KeyParams(),
	}
	return self.renderSuccess(result)
}

func (self *loginHandler) handleError(err error) error {
	var odinError *odinerrors.Error
	if !errors.As(err, &odinError) {
		return err
	}
	switch odinError.Tag() {
	case odinerrors.Unauthorized:
		self.ctx.Status(http.StatusUnauthorized)
	case odinerrors.Domain:
		self.ctx.Status(http.StatusBadRequest)
	default:
		return err
	}
	return self.renderError(err)
}

func (self *loginHandler) renderSuccess(result *LoginResult) error {
	return self.ctx.JSON(loginResponse{
		Token:              result.Token,
		EncryptedMasterKey: result.EncryptedMasterKey,
		KeyParams: &keyParamsResponse{
			Algorithm:   result.KeyParams.Algorithm(),
			Iterations:  result.KeyParams.Iterations(),
			Memory:      result.KeyParams.Memory(),
			Parallelism: result.KeyParams.Parallelism(),
			Salt:        result.KeyParams.Salt(),
		},
	})
}

func (self *loginHandler) renderError(err error) error {
	var odinError *odinerrors.Error
	errors.As(err, &odinError)
	return self.ctx.JSON(loginResponse{Error: odinError.ExternalError()})
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
