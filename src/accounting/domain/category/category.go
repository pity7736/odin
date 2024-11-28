package categorymodel

import (
	"raiseexception.dev/odin/src/accounting/domain/constants"
)

type Category struct {
	name   string
	id     string
	t      constants.CategoryType
	userID string
}

func New(id, name string, t constants.CategoryType, userID string) *Category {
	return &Category{
		name:   name,
		id:     id,
		t:      t,
		userID: userID,
	}
}

func (self *Category) Name() string {
	return self.name
}

func (self *Category) ID() string {
	return self.id
}

func (self *Category) Type() constants.CategoryType {
	return self.t
}

func (self *Category) UserID() string {
	return self.userID
}
