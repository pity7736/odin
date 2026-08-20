package inmemory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders/userbuilder"
)

func TestInMemoryUserRepository(t *testing.T) {
	t.Run("GetByEmail returns a NotFound error when the user is absent", func(t *testing.T) {
		repository := inmemory.NewInMemoryUserRepository()

		user, err := repository.GetByEmail(context.TODO(), "missing@example.com")

		assert.Nil(t, user)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.NotFound, odinError.Tag())
	})

	t.Run("GetByEmail returns the user when present", func(t *testing.T) {
		repository := inmemory.NewInMemoryUserRepository()
		stored := userbuilder.New().WithEmail("present@example.com").Create(repository)

		user, err := repository.GetByEmail(context.TODO(), "present@example.com")

		assert.Nil(t, err)
		assert.Equal(t, stored, user)
	})

	t.Run("Exists returns false when the user is absent", func(t *testing.T) {
		repository := inmemory.NewInMemoryUserRepository()

		exists, err := repository.Exists(context.TODO(), "missing@example.com")

		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists returns true after the user is added", func(t *testing.T) {
		repository := inmemory.NewInMemoryUserRepository()
		userbuilder.New().WithEmail("present@example.com").Create(repository)

		exists, err := repository.Exists(context.TODO(), "present@example.com")

		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
