package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/tests/testutils"
)

const tokenLength = 50

func TestSessionShould(t *testing.T) {
	t.Run("generate a token of expected length", func(t *testing.T) {
		session, err := sessionmodel.New("user-id", sessionmodel.DefaultTTL)
		assert.Nil(t, err)
		assert.Len(t, session.Token(), tokenLength)
	})
	t.Run("set expires at approximately now plus ttl", func(t *testing.T) {
		ttl := sessionmodel.DefaultTTL
		expectedExpiry := time.Now().Add(ttl)
		session, _ := sessionmodel.New("user-id", ttl)
		assert.True(t, testutils.IsTimeClose(expectedExpiry, session.ExpiresAt()))
	})
	t.Run("not be expired for a new session", func(t *testing.T) {
		session, _ := sessionmodel.New("user-id", sessionmodel.DefaultTTL)
		assert.False(t, session.IsExpired())
	})
	t.Run("be expired when past expiry", func(t *testing.T) {
		pastExpiry := time.Now().Add(-time.Hour)
		session := sessionmodel.NewFromRepository("token", "user-id", time.Now().Add(-2*time.Hour), pastExpiry)
		assert.True(t, session.IsExpired())
	})
	t.Run("store created at time", func(t *testing.T) {
		before := time.Now()
		session, _ := sessionmodel.New("user-id", sessionmodel.DefaultTTL)
		assert.True(t, testutils.IsTimeClose(before, session.CreatedAt()))
	})
	t.Run("extend resets expires at to now plus ttl", func(t *testing.T) {
		session, _ := sessionmodel.New("user-id", sessionmodel.DefaultTTL)
		ttl := 24 * time.Hour
		expectedExpiry := time.Now().Add(ttl)
		session.Extend(ttl)
		assert.True(t, testutils.IsTimeClose(expectedExpiry, session.ExpiresAt()))
	})
}
