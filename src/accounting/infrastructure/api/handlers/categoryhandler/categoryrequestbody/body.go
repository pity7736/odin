package categoryrequestbody

import (
	"errors"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/user"
)

type CategoryRequestBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
	User string `json:"user"`
}

func (c CategoryRequestBody) CreateCategoryCreatorCommand() (*categorycommand.CategoryCreatorCommand, error) {
	if c.Name != "" {
		categoryType, _ := constants.NewFromString(c.Type)
		command := categorycommand.New(
			c.Name,
			categoryType,
			user.New(c.User),
		)
		return &command, nil
	}
	return nil, errors.New("bad request")
}
