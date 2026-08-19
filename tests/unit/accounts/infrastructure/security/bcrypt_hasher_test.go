package security_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
)

func TestBcryptHasherShould(t *testing.T) {
	t.Run("produce a digest that matches the original value", func(t *testing.T) {
		hasher := bcrypthasher.New()
		authHash := "client-derived-auth-hash"
		digest, err := hasher.Hash(authHash)
		assert.Nil(t, err)
		assert.NotEmpty(t, digest)
		assert.True(t, hasher.Compare(digest, authHash))
	})
	t.Run("produce a digest that does not match a different value", func(t *testing.T) {
		hasher := bcrypthasher.New()
		digest, err := hasher.Hash("client-derived-auth-hash")
		assert.Nil(t, err)
		assert.False(t, hasher.Compare(digest, "other"))
	})
	t.Run("return an error when the input exceeds the supported length", func(t *testing.T) {
		hasher := bcrypthasher.New()
		digest, err := hasher.Hash(strings.Repeat("a", 73))
		assert.NotNil(t, err)
		assert.Empty(t, digest)
	})
}
