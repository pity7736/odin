package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

func TestUserShould(t *testing.T) {
	t.Run("check password successfully", func(t *testing.T) {
		password := "secure password"
		user, err := usermodel.New("test@example.com", password)
		assert.Nil(t, err)
		assert.True(t, user.CheckPassword(password))
	})
	t.Run("fail check password with wrong password", func(t *testing.T) {
		user, _ := usermodel.New("test@example.com", "correct password")
		assert.False(t, user.CheckPassword("wrong password"))
	})
	t.Run("not store plain text password", func(t *testing.T) {
		password := "plain text password"
		user, _ := usermodel.New("test@example.com", password)
		assert.NotEqual(t, password, user.HashedPassword())
	})
	t.Run("generate a unique ID", func(t *testing.T) {
		user, _ := usermodel.New("test@example.com", "password")
		assert.NotEmpty(t, user.ID())
	})
}
