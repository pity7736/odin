package categorymodel

import (
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
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

func (self *Category) ValidateOwnership(requestContext *requestcontext.RequestContext) error {
	if self.UserID() != requestContext.UserID() {
		return odinerrors.NewErrorBuilder("categoría no pertenece a usuario logueado").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("la categoría no pertenece al usuario logueado").
			Build()
	}
	return nil
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

func (self *Category) IsIncome() bool {
	return self.t == constants.INCOME
}
