package create_category__test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounting/application/commands/categorycommand"
	"raiseexception.dev/odin/src/accounting/application/use_cases/categorycreator"
	"raiseexception.dev/odin/src/accounting/domain/constants"
	"raiseexception.dev/odin/src/shared/domain/user"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type setup struct {
	repository   *mocks.MockCategoryRepository
	command      categorycommand.CategoryCreatorCommand
	categoryName string
	categoryType constants.CategoryType
	user         *user.User
}

func newSetup(t *testing.T) setup {
	repository := mocks.NewMockCategoryRepository(t)
	categoryName := "test"
	user := getUser()
	categoryType := constants.EXPENSE
	command := categorycommand.New(categoryName, categoryType, user)
	return setup{
		repository:   repository,
		command:      command,
		categoryName: categoryName,
		categoryType: categoryType,
		user:         user,
	}
}

func Test(t *testing.T) {

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		setup := newSetup(t)
		id, _ := uuid.NewV7()
		patches := gomonkey.ApplyFuncReturn(uuid.NewV7, id, nil)
		defer patches.Reset()
		user := user.New("test@raiseexception.dev")
		setup.repository.On("Add", ctx, mock.Anything).Return(nil)
		categoryCreator := categorycreator.New(setup.command, setup.repository)
		category, err := categoryCreator.Create(ctx)

		assert.Nil(t, err)
		assert.Equal(t, id.String(), category.ID())
		assert.Equal(t, setup.categoryName, category.Name())
		assert.Equal(t, setup.categoryType, category.Type())
		assert.Equal(t, *user, category.User())
		setup.repository.AssertCalled(t, "Add", ctx, category)
	})

	t.Run("when repository fails", func(t *testing.T) {
		setup := newSetup(t)
		savingError := fmt.Errorf("error saving the category %s", setup.categoryName)
		setup.repository.On("Add", ctx, mock.Anything).Return(savingError)
		categoryCreator := categorycreator.New(setup.command, setup.repository)
		category, err := categoryCreator.Create(ctx)

		assert.Nil(t, category)
		assert.Equal(t, savingError, err)
	})

	t.Run("when id generation fails", func(t *testing.T) {
		setup := newSetup(t)
		patches := gomonkey.ApplyFuncReturn(uuid.NewV7, nil, errors.New("some error"))
		defer patches.Reset()
		categoryCreator := categorycreator.New(setup.command, setup.repository)
		category, err := categoryCreator.Create(ctx)

		assert.Nil(t, category)
		assert.Equal(t, errors.New("error generating uuid"), err)
		setup.repository.AssertNotCalled(t, "Add")
	})
}

func getUser() *user.User {
	return user.New("test@raiseexception.dev")
}
