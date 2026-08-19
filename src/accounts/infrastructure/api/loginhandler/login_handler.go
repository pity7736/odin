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

type LoginHandler interface {
	HandleResponse(result *LoginResult) error
	HandleBadRequest(err error) error
	ContentType() string
}

type loginHandler struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	authHasher        authhasher.AuthHasher
	handler           LoginHandler
}

func New(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, authHasher authhasher.AuthHasher, handler LoginHandler) *loginHandler {
	return &loginHandler{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		authHasher:        authHasher,
		handler:           handler,
	}
}

func (self *loginHandler) Login(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", self.handler.ContentType())
	var body LoginBody
	if err := self.validateRequestBody(ctx, &body); err != nil {
		ctx.Status(http.StatusBadRequest)
		return self.handler.HandleBadRequest(err)
	}
	return self.login(ctx, &body)
}

func (self *loginHandler) validateRequestBody(ctx *fiber.Ctx, body *LoginBody) error {
	if err := ctx.BodyParser(body); err != nil {
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

func (self *loginHandler) login(ctx *fiber.Ctx, body *LoginBody) error {
	starter := sessionstarter.New(
		strings.Clone(body.Email),
		strings.Clone(body.AuthHash),
		self.userRepository,
		self.sessionRepository,
		self.authHasher,
	)
	session, user, err := starter.Start(ctx.Context())
	if err != nil {
		return self.handleError(ctx, err)
	}
	ctx.Status(http.StatusCreated)
	return self.handler.HandleResponse(&LoginResult{
		Token:              session.Token(),
		EncryptedMasterKey: user.EncryptedMasterKey(),
		KeyParams:          user.KeyParams(),
	})
}

func (self *loginHandler) handleError(ctx *fiber.Ctx, err error) error {
	var odinError *odinerrors.Error
	if !errors.As(err, &odinError) {
		return err
	}
	switch odinError.Tag() {
	case odinerrors.Unauthorized:
		ctx.Status(http.StatusUnauthorized)
	case odinerrors.Domain:
		ctx.Status(http.StatusBadRequest)
	default:
		return err
	}
	return self.handler.HandleBadRequest(err)
}
