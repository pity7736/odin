package categorybuilder

import (
	"context"

	"github.com/google/uuid"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/shared/domain/user"
	"raiseexception.dev/odin/tests/builders/userbuilder"
)

type builder struct {
	name string
	id   string
	t    constants.CategoryType
	user *user.User
}

func New() *builder {
	id, _ := uuid.NewV7()
	return &builder{
		name: "test",
		id:   id.String(),
		t:    constants.EXPENSE,
	}
}

func (b *builder) WithDefaultUser() *builder {
	b.user = userbuilder.New().Build()
	return b
}

func (b *builder) Build() *category.Category {
	return category.New(
		b.id,
		b.name,
		b.t,
		b.user,
	)
}

func (b *builder) Create(repository repositories.CategoryRepository) *category.Category {
	command := categorycommand.New(b.name, b.t, b.user)
	categoryCreator := categorycreator.New(command, repository)
	category, _ := categoryCreator.Create(context.TODO())
	return category
}
