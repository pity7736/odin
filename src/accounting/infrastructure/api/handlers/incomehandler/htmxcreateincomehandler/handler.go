package htmxcreateincomehandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"raiseexception.dev/odin/src/accounting/application/use_cases/incomecreator"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type HTMXCreateIncomeHandler struct {
	accountRepository repositories.AccountRepository
	factory           accountingrepositoryfactory.RepositoryFactory
}

func New(factory accountingrepositoryfactory.RepositoryFactory) *HTMXCreateIncomeHandler {
	return &HTMXCreateIncomeHandler{
		accountRepository: factory.GetAccountRepository(),
		factory:           factory,
	}
}

func (self *HTMXCreateIncomeHandler) Handle(c *fiber.Ctx) error {
	body, err := newFromCtx(c)
	if err != nil {
		c.Render("create_account_error", err, "")
		return err
	}
	command, _ := body.toCommand()
	incomeCreator := incomecreator.New(self.factory, command)
	requestContext := c.Locals(requestcontext.Key).(*requestcontext.RequestContext)
	ctx := context.WithValue(c.Context(), requestcontext.Key, requestContext)
	income, err := incomeCreator.Create(ctx)
	if err != nil {
		c.Render("create_account_error", err, "")
		return err
	}
	c.Render("income_created", income, "")
	return err
}
