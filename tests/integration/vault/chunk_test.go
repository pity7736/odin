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

func TestReadChunkIntegrationShould(t *testing.T) {
	t.Run("return a stored chunk to its authenticated owner exactly as stored", func(t *testing.T) {
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
		content := "nonce-ciphertext-tag-base64"
		body := fmt.Sprintf(`{"id": "%s", "content": "%s"}`, id, content)
		createBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/chunks").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON).
			WithUser(owner)
		createResponse := testutils.GetJSONResponseFromRequestBuilder(application, createBuilder)
		defer func() { _ = createResponse.Body.Close() }()
		assert.Equal(t, http.StatusCreated, createResponse.StatusCode)
		var responseData map[string]any
		readBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath("/api/v1/chunks/" + id).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(application, readBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, id, responseData["id"])
		assert.Equal(t, content, responseData["content"])
	})

	t.Run("hide a chunk from a user who is not its owner", func(t *testing.T) {
		userRepository := accountsinmemory.NewInMemoryUserRepository()
		sessionRepository := accountsinmemory.NewInMemorySessionRepository()
		chunkRepository := inmemory.NewInMemoryChunkRepository()
		application := apptest.New().
			WithUserRepository(userRepository).
			WithSessionRepository(sessionRepository).
			WithChunkRepository(chunkRepository).
			Build()
		firstOwner := userbuilder.New().WithEmail("first@example.com").Create(userRepository)
		otherUser := userbuilder.New().WithEmail("other@example.com").Create(userRepository)
		id := uuid.NewString()
		body := fmt.Sprintf(`{"id": "%s", "content": "nonce-ciphertext-tag-base64"}`, id)
		createBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithPath("/api/v1/chunks").
			WithPayload(body).
			WithContentType(fiber.MIMEApplicationJSON).
			WithUser(firstOwner)
		createResponse := testutils.GetJSONResponseFromRequestBuilder(application, createBuilder)
		defer func() { _ = createResponse.Body.Close() }()
		assert.Equal(t, http.StatusCreated, createResponse.StatusCode)
		var responseData map[string]any
		readBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath("/api/v1/chunks/" + id).
			WithResponseData(&responseData).
			WithUser(otherUser)
		response := testutils.GetJSONResponseFromRequestBuilder(application, readBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, "El elemento no existe", responseData["error"])
	})
}
