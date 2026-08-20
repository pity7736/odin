package vault_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	accountsinmemory "raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/vault/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/apptest"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func TestCreateChunkIntegrationShould(t *testing.T) {
	t.Run("store a chunk for the authenticated owner and acknowledge it", func(t *testing.T) {
		userRepository := accountsinmemory.NewInMemoryUserRepository()
		sessionRepository := accountsinmemory.NewInMemorySessionRepository()
		chunkRepository := inmemory.NewInMemoryChunkRepository()
		application := apptest.New().
			WithUserRepository(userRepository).
			WithSessionRepository(sessionRepository).
			WithChunkRepository(chunkRepository).
			Build()
		owner := userbuilder.New().WithEmail("owner@example.com").Create(userRepository)
		id := uuid.NewString()
		body := fmt.Sprintf(`{"id": "%s", "content": "nonce-ciphertext-tag-base64"}`, id)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/chunks").
			WithPayload(body).
			WithResponseData(&responseData).
			WithContentType(fiber.MIMEApplicationJSON).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, id, responseData["id"])
		exists, err := chunkRepository.Exists(context.TODO(), owner.ID(), id)
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
