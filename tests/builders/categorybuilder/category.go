package categorybuilder

import (
	"context"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"

	"github.com/google/uuid"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
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

func (b *builder) Build() *category.Category {
	return category.New(
		b.id,
		b.name,
		b.t,
		b.userID,
	)
}

func (b *builder) Create(repository repositories.CategoryRepository) *category.Category {
	command := categorycommand.New(b.name, b.t, b.userID)
	categoryCreator := categorycreator.New(command, repository)
	category, _ := categoryCreator.Create(context.TODO())
	return category
}
