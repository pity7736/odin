package restloginhandler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type restLoginHandler struct {
	ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) loginhandler.LoginHandler {
	return &restLoginHandler{ctx: ctx}
}

func (self *restLoginHandler) HandleResponse(result *loginhandler.LoginResult) error {
	return self.ctx.JSON(response{
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

func (self *restLoginHandler) HandleBadRequest(err error) error {
	var odinError *odinerrors.Error
	errors.As(err, &odinError)
	return self.ctx.JSON(response{Error: odinError.ExternalError()})
}

func (self *restLoginHandler) ContentType() string {
	return fiber.MIMEApplicationJSON
}

type response struct {
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
