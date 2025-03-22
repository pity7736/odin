package requestcontext_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

func TestRequestContextShould(t *testing.T) {
	t.Run("return error when user id is empty", func(t *testing.T) {
		ctx, err := requestcontext.New("")

		assert.Equal(t, "user id cannot be empty", err.Error())
		assert.Nil(t, ctx)
	})

	t.Run("return anonymous when request id is empty", func(t *testing.T) {
		ctx := requestcontext.NewAnonymous()

		assert.True(t, ctx.IsAnonymous())
	})

	t.Run("return not anonymous when data is valid", func(t *testing.T) {
		userID := "12345"
		ctx, err := requestcontext.New(userID)

		assert.Nil(t, err)
		assert.Equal(t, userID, ctx.UserID())
		assert.False(t, ctx.IsAnonymous())
	})
}
