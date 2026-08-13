package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

func TestUserShould(t *testing.T) {
	t.Run("generate a unique ID", func(t *testing.T) {
		user, err := usermodel.New("test@example.com", "hashed-password")
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID())
	})
	t.Run("store email and hashed password", func(t *testing.T) {
		user, _ := usermodel.New("test@example.com", "hashed-password")
		assert.Equal(t, "test@example.com", user.Email())
		assert.Equal(t, "hashed-password", user.HashedPassword())
	})
	t.Run("generate different IDs for different users", func(t *testing.T) {
		user1, _ := usermodel.New("a@example.com", "hash1")
		user2, _ := usermodel.New("b@example.com", "hash2")
		assert.NotEqual(t, user1.ID(), user2.ID())
	})
}
