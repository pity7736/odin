package builders

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
	"raiseexception.dev/odin/tests/builders/userbuilder"
)

type FiberContextBuilder struct {
	body           []byte
	contentType    string
	method         string
	ctx            *fiber.Ctx
	app            *fiber.App
	user           *usermodel.User
	requestContext *requestcontext.RequestContext
}

func NewFiberContextBuilder() *FiberContextBuilder {
	engine := html.New(
		"/Users/julian.cortes/development/odin/src/shared/infrastructure/templates",
		".gohtml",
	)
	return &FiberContextBuilder{
		method: fiber.MethodGet,
		app: fiber.New(fiber.Config{
			Views:       engine,
			ViewsLayout: "base",
		}),
		user:           userbuilder.New().Build(),
		requestContext: requestcontext.NewAnonymous(),
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

func (self *FiberContextBuilder) WithAnonymousRequest() *FiberContextBuilder {
	self.user = nil
	return self
}

func (self *FiberContextBuilder) WithoutRequestContext() *FiberContextBuilder {
	self.WithAnonymousRequest()
	self.requestContext = nil
	return self
}

func (self *FiberContextBuilder) Build() *fiber.Ctx {
	self.ctx = self.app.AcquireCtx(&fasthttp.RequestCtx{})
	self.ctx.Method(self.method)
	if self.body != nil {
		self.ctx.Request().SetBody(self.body)
		self.ctx.Request().Header.SetContentLength(len(self.body))
	}
	if self.contentType != "" {
		self.ctx.Request().Header.SetContentType(self.contentType)
	}
	if self.user != nil {
		self.ctx.Locals("userID", self.user.ID())
		self.requestContext, _ = requestcontext.New(self.user.ID())
	}
	if self.requestContext != nil {
		self.ctx.Locals(requestcontext.Key, self.requestContext)
	}
	return self.ctx
}

func (self *FiberContextBuilder) Release() {
	self.app.ReleaseCtx(self.ctx)
}

func (self *FiberContextBuilder) User() *usermodel.User {
	return self.user
}
