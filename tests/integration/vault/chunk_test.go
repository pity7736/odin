package vault_test

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	accountsinmemory "raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/app"
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

func storeChunk(t *testing.T, application app.Application, userRepository *accountsinmemory.InMemoryUserRepository, sessionRepository *accountsinmemory.InMemorySessionRepository, owner *usermodel.User, id, content string) {
	t.Helper()
	body := fmt.Sprintf(`{"id": "%s", "content": "%s"}`, id, content)
	requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
		WithPath("/api/v1/chunks").
		WithPayload(body).
		WithContentType(fiber.MIMEApplicationJSON).
		WithUser(owner)
	response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
	defer func() { _ = response.Body.Close() }()
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestListChunkIntegrationShould(t *testing.T) {
	t.Run("return all the owner's chunks ordered newest first", func(t *testing.T) {
		userRepository := accountsinmemory.NewInMemoryUserRepository()
		sessionRepository := accountsinmemory.NewInMemorySessionRepository()
		chunkRepository := inmemory.NewInMemoryChunkRepository()
		application := apptest.New().
			WithUserRepository(userRepository).
			WithSessionRepository(sessionRepository).
			WithChunkRepository(chunkRepository).
			Build()
		owner := userbuilder.New().WithEmail("owner@example.com").Create(userRepository)
		firstID := uuid.NewString()
		secondID := uuid.NewString()
		thirdID := uuid.NewString()
		storeChunk(t, application, userRepository, sessionRepository, owner, firstID, "first-content")
		storeChunk(t, application, userRepository, sessionRepository, owner, secondID, "second-content")
		storeChunk(t, application, userRepository, sessionRepository, owner, thirdID, "third-content")
		expectedOrder := []string{firstID, secondID, thirdID}
		sort.Sort(sort.Reverse(sort.StringSlice(expectedOrder)))
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath("/api/v1/chunks").
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		chunks := responseData["chunks"].([]any)
		assert.Len(t, chunks, 3)
		assert.Equal(t, expectedOrder[0], chunks[0].(map[string]any)["id"])
		assert.Equal(t, expectedOrder[1], chunks[1].(map[string]any)["id"])
		assert.Equal(t, expectedOrder[2], chunks[2].(map[string]any)["id"])
	})

	t.Run("return only the requesting owner's chunks", func(t *testing.T) {
		userRepository := accountsinmemory.NewInMemoryUserRepository()
		sessionRepository := accountsinmemory.NewInMemorySessionRepository()
		chunkRepository := inmemory.NewInMemoryChunkRepository()
		application := apptest.New().
			WithUserRepository(userRepository).
			WithSessionRepository(sessionRepository).
			WithChunkRepository(chunkRepository).
			Build()
		owner := userbuilder.New().WithEmail("owner@example.com").Create(userRepository)
		otherUser := userbuilder.New().WithEmail("other@example.com").Create(userRepository)
		ownerID := uuid.NewString()
		otherID := uuid.NewString()
		storeChunk(t, application, userRepository, sessionRepository, owner, ownerID, "owner-content")
		storeChunk(t, application, userRepository, sessionRepository, otherUser, otherID, "other-content")
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath("/api/v1/chunks").
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		chunks := responseData["chunks"].([]any)
		assert.Len(t, chunks, 1)
		assert.Equal(t, ownerID, chunks[0].(map[string]any)["id"])
		assert.Equal(t, "owner-content", chunks[0].(map[string]any)["content"])
	})

	t.Run("return an empty collection for an owner with no chunks", func(t *testing.T) {
		userRepository := accountsinmemory.NewInMemoryUserRepository()
		sessionRepository := accountsinmemory.NewInMemorySessionRepository()
		chunkRepository := inmemory.NewInMemoryChunkRepository()
		application := apptest.New().
			WithUserRepository(userRepository).
			WithSessionRepository(sessionRepository).
			WithChunkRepository(chunkRepository).
			Build()
		owner := userbuilder.New().WithEmail("owner@example.com").Create(userRepository)
		var responseData map[string]any
		requestBuilder := builders.NewRequestBuilder(userRepository, sessionRepository).
			WithMethod("GET").
			WithPath("/api/v1/chunks").
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		chunks := responseData["chunks"].([]any)
		assert.Empty(t, chunks)
	})
}
