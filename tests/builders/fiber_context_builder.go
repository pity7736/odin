package builders

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/tests/builders/userbuilder"
)

type FiberContextBuilder struct {
	body        []byte
	contentType string
	method      string
	ctx         *fiber.Ctx
	app         *fiber.App
	user        *usermodel.User
}

func NewFiberContextBuilder() *FiberContextBuilder {
	return &FiberContextBuilder{
		method: fiber.MethodGet,
		app:    fiber.New(),
		user:   userbuilder.New().Build(),
	}
}

func (self *FiberContextBuilder) WithBody(body []byte) *FiberContextBuilder {
	self.body = body
	return self
}

func (self *FiberContextBuilder) WithContentType(contentType string) *FiberContextBuilder {
	self.contentType = contentType
	return self
}

func (self *FiberContextBuilder) WithMethod(method string) *FiberContextBuilder {
	self.method = method
	return self
}

func (self *FiberContextBuilder) Build() *fiber.Ctx {
	self.ctx = self.app.AcquireCtx(&fasthttp.RequestCtx{})
	if self.body != nil {
		self.ctx.Request().SetBody(self.body)
		self.ctx.Request().Header.SetContentLength(len(self.body))
	}
	if self.contentType != "" {
		self.ctx.Request().Header.SetContentType(self.contentType)
	}
	if self.user != nil {
		self.ctx.Locals("userID", self.user.ID())
	}
	return self.ctx
}

func (self *FiberContextBuilder) Release() {
	self.app.ReleaseCtx(self.ctx)
}

func (self *FiberContextBuilder) User() *usermodel.User {
	return self.user
}
