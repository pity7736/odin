package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

func TestUserShould(t *testing.T) {
	t.Run("generate a unique ID", func(t *testing.T) {
		params, _ := keyparams.New("argon2id", 3, 65536, 4, "salt")
		user, err := usermodel.New("test@example.com", "auth-hash-digest", "encrypted-master-key", params)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID())
	})
	t.Run("store email and auth hash digest", func(t *testing.T) {
		params, _ := keyparams.New("argon2id", 3, 65536, 4, "salt")
		user, _ := usermodel.New("test@example.com", "auth-hash-digest", "encrypted-master-key", params)
		assert.Equal(t, "test@example.com", user.Email())
		assert.Equal(t, "auth-hash-digest", user.AuthHashDigest())
	})
	t.Run("generate different IDs for different users", func(t *testing.T) {
		params, _ := keyparams.New("argon2id", 3, 65536, 4, "salt")
		user1, _ := usermodel.New("a@example.com", "hash1", "emk1", params)
		user2, _ := usermodel.New("b@example.com", "hash2", "emk2", params)
		assert.NotEqual(t, user1.ID(), user2.ID())
	})
	t.Run("store encrypted master key and key params", func(t *testing.T) {
		params, _ := keyparams.New("argon2id", 3, 65536, 4, "salt")
		user, _ := usermodel.New("test@example.com", "auth-hash-digest", "encrypted-master-key", params)
		assert.Equal(t, "encrypted-master-key", user.EncryptedMasterKey())
		assert.Equal(t, params, user.KeyParams())
	})
}
