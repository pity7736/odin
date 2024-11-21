package categorycommand

import (
	"raiseexception.dev/odin/src/accounting/domain/constants"
)

type CategoryCreatorCommand struct {
	name         string
	categoryType constants.CategoryType
	userID       string
}

func New(name string, categoryType constants.CategoryType, userID string) CategoryCreatorCommand {
	return CategoryCreatorCommand{name: name, categoryType: categoryType, userID: userID}
}

func (c CategoryCreatorCommand) Name() string {
	return c.name
}

func (c CategoryCreatorCommand) Type() constants.CategoryType {
	return c.categoryType
}

func (c CategoryCreatorCommand) UserID() string {
	return c.userID
}
