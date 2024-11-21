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

func (self CategoryCreatorCommand) Name() string {
	return self.name
}

func (self CategoryCreatorCommand) Type() constants.CategoryType {
	return self.categoryType
}

func (self CategoryCreatorCommand) UserID() string {
	return self.userID
}
