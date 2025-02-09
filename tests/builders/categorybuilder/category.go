package categorybuilder

import (
	"context"

	"github.com/google/uuid"

	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type builder struct {
	name   string
	id     string
	t      constants.CategoryType
	user   *usermodel.User
	userID string
}

func New() *builder {
	id, _ := uuid.NewV7()
	userID, _ := uuid.NewV7()
	return &builder{
		name:   "test",
		id:     id.String(),
		t:      constants.EXPENSE,
		userID: userID.String(),
	}
}

func (self *builder) Build() *categorymodel.Category {
	return categorymodel.New(
		self.id,
		self.name,
		self.t,
		self.userID,
	)
}

func (self *builder) Create(repository repositories.CategoryRepository) *categorymodel.Category {
	command := categorycommand.New(self.name, self.t, self.userID)
	categoryCreator := categorycreator.New(command, repository)
	category, _ := categoryCreator.Create(context.TODO())
	return category
}
