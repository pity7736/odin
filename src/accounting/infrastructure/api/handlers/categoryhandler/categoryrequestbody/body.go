package categoryrequestbody

import (
	"errors"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/constants"
)

type CategoryRequestBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c CategoryRequestBody) CreateCategoryCreatorCommand(userID string) (*categorycommand.CategoryCreatorCommand, error) {
	if c.Name != "" && c.Type != "" {
		categoryType, _ := constants.NewFromString(c.Type)
		command := categorycommand.New(
			c.Name,
			categoryType,
			userID,
		)
		return &command, nil
	}
	return nil, errors.New("bad request")
}
