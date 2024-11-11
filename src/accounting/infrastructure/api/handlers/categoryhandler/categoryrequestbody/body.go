package categoryrequestbody

import (
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/user"
)

type CategoryRequestBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
	User string `json:"user"`
}

func (c CategoryRequestBody) CreateCategoryCreatorCommand() categorycommand.CategoryCreatorCommand {
	categoryType, _ := constants.NewFromString(c.Type)
	command := categorycommand.New(
		c.Name,
		categoryType,
		user.New(c.User),
	)
	return command
}
