package loginhandler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounts/application/passwordhasher"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type LoginHandler interface {
	HandleResponse(session *sessionmodel.Session) error
	HandleBadRequest(err error) error
	ContentType() string
}

type loginHandler struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
	passwordHasher    passwordhasher.PasswordHasher
	handler           LoginHandler
}

func New(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, passwordHasher passwordhasher.PasswordHasher, handler LoginHandler) *loginHandler {
	return &loginHandler{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		passwordHasher:    passwordHasher,
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
			WithTag(odinerrors.DOMAIN).
			Build()
	}
	if body.Email == "" {
		return odinerrors.NewErrorBuilder("email is required").
			WithExternalMessage("El correo es obligatorio").
			WithTag(odinerrors.DOMAIN).
			Build()
	}
	if body.Password == "" {
		return odinerrors.NewErrorBuilder("password is required").
			WithExternalMessage("La contraseña es obligatoria").
			WithTag(odinerrors.DOMAIN).
			Build()
	}
	return nil
}

func (self *loginHandler) login(ctx *fiber.Ctx, body *LoginBody) error {
	starter := sessionstarter.New(
		strings.Clone(body.Email),
		strings.Clone(body.Password),
		self.userRepository,
		self.sessionRepository,
		self.passwordHasher,
	)
	session, err := starter.Start(ctx.Context())
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return self.handler.HandleBadRequest(err)
	}
	ctx.Status(http.StatusCreated)
	return self.handler.HandleResponse(session)
}
