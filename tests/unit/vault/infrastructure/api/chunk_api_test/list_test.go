package chunk_api_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func (self chunkHandlerTestContext) listRequestBuilder() *builders.RequestBuilder {
	return builders.NewRequestBuilder(self.userRepository, self.sessionRepository).
		WithMethod("GET").
		WithPath("/api/v1/chunks")
}

func TestListChunksRestShould(t *testing.T) {
	t.Run("return every chunk of the authenticated owner as stored", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		firstID := uuid.NewString()
		secondID := uuid.NewString()
		firstChunk, _ := chunkmodel.New(firstID, owner.ID(), "first-content")
		secondChunk, _ := chunkmodel.New(secondID, owner.ID(), "second-content")
		storedChunks := []*chunkmodel.EncryptedChunk{firstChunk, secondChunk}
		testContext.chunkRepository.EXPECT().GetAll(mock.Anything, owner.ID()).Return(storedChunks, nil)
		var responseData map[string]any
		requestBuilder := testContext.listRequestBuilder().
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		chunks := responseData["chunks"].([]any)
		assert.Len(t, chunks, 2)
		first := chunks[0].(map[string]any)
		assert.Equal(t, firstID, first["id"])
		assert.Equal(t, "first-content", first["content"])
		second := chunks[1].(map[string]any)
		assert.Equal(t, secondID, second["id"])
		assert.Equal(t, "second-content", second["content"])
	})

	t.Run("return an empty collection when the owner has no chunks", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		testContext.chunkRepository.EXPECT().GetAll(mock.Anything, owner.ID()).Return([]*chunkmodel.EncryptedChunk{}, nil)
		var responseData map[string]any
		requestBuilder := testContext.listRequestBuilder().
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		chunks := responseData["chunks"].([]any)
		assert.Empty(t, chunks)
	})

	t.Run("reject an unauthenticated request", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		requestBuilder := testContext.listRequestBuilder().
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		testContext.chunkRepository.AssertNotCalled(t, "GetAll", mock.Anything, mock.Anything)
	})
}
