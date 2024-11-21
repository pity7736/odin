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

func (self CategoryRequestBody) CreateCategoryCreatorCommand(userID string) (*categorycommand.CategoryCreatorCommand, error) {
	if self.Name != "" && self.Type != "" {
		categoryType, _ := constants.NewFromString(self.Type)
		command := categorycommand.New(
			self.Name,
			categoryType,
			userID,
		)
		return &command, nil
	}
	return nil, errors.New("bad request")
}
