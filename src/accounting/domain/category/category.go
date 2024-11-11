package category

import (
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/user"
)

type Category struct {
	name string
	id   string
	t    constants.CategoryType
	user *user.User
}

func New(id, name string, t constants.CategoryType, user *user.User) *Category {
	return &Category{
		name: name,
		id:   id,
		t:    t,
		user: user,
	}
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) ID() string {
	return c.id
}

func (c *Category) Type() constants.CategoryType {
	return c.t
}

func (c *Category) User() *user.User {
	return c.user
}
