package categorycreator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
)

type CategoryCreator struct {
	command    categorycommand.CategoryCreatorCommand
	repository repositories.CategoryRepository
}

func New(command categorycommand.CategoryCreatorCommand, repository repositories.CategoryRepository) CategoryCreator {
	return CategoryCreator{command: command, repository: repository}
}

func (self CategoryCreator) Create(ctx context.Context) (*categorymodel.Category, error) {
	id, uuidError := uuid.NewV7()
	if uuidError != nil {
		return nil, errors.New("error generating uuid")
	}
	category := categorymodel.New(
		id.String(),
		self.command.Name(),
		self.command.Type(),
		self.command.UserID(),
	)
	err := self.repository.Add(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}
