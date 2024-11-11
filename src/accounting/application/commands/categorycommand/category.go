package categorycommand

import (
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/user"
)

type CategoryCreatorCommand struct {
	name         string
	categoryType constants.CategoryType
	user         *user.User
}

func New(name string, categoryType constants.CategoryType, user *user.User) CategoryCreatorCommand {
	return CategoryCreatorCommand{name: name, categoryType: categoryType, user: user}
}

func (c CategoryCreatorCommand) Name() string {
	return c.name
}

func (c CategoryCreatorCommand) Type() constants.CategoryType {
	return c.categoryType
}

func (c CategoryCreatorCommand) User() *user.User {
	return c.user
}
