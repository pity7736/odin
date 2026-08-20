package chunk_api_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func (self chunkHandlerTestContext) readRequestBuilder(id string) *builders.RequestBuilder {
	return builders.NewRequestBuilder(self.userRepository, self.sessionRepository).
		WithMethod("GET").
		WithPath("/api/v1/chunks/" + id)
}

func TestReadChunkRestShould(t *testing.T) {
	t.Run("return a stored chunk to its authenticated owner", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		storedChunk, _ := chunkmodel.New(id, owner.ID(), "encrypted-content")
		testContext.chunkRepository.EXPECT().Get(mock.Anything, owner.ID(), id).Return(storedChunk, nil)
		var responseData map[string]any
		requestBuilder := testContext.readRequestBuilder(id).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, id, responseData["id"])
		assert.Equal(t, "encrypted-content", responseData["content"])
	})

	t.Run("return not found when the owner has no chunk with that id", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		notFoundError := odinerrors.NewErrorBuilder("chunk not found").
			WithExternalMessage("El elemento no existe").
			WithTag(odinerrors.NotFound).
			Build()
		testContext.chunkRepository.EXPECT().Get(mock.Anything, owner.ID(), id).Return(nil, notFoundError)
		var responseData map[string]any
		requestBuilder := testContext.readRequestBuilder(id).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, "El elemento no existe", responseData["error"])
	})

	t.Run("return not found for a chunk that belongs to another user", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		notFoundError := odinerrors.NewErrorBuilder("chunk not found").
			WithExternalMessage("El elemento no existe").
			WithTag(odinerrors.NotFound).
			Build()
		testContext.chunkRepository.EXPECT().Get(mock.Anything, owner.ID(), id).Return(nil, notFoundError)
		var responseData map[string]any
		requestBuilder := testContext.readRequestBuilder(id).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, "El elemento no existe", responseData["error"])
	})

	t.Run("pass a malformed id straight to the lookup and return not found", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		notFoundError := odinerrors.NewErrorBuilder("chunk not found").
			WithExternalMessage("El elemento no existe").
			WithTag(odinerrors.NotFound).
			Build()
		testContext.chunkRepository.EXPECT().Get(mock.Anything, owner.ID(), "not-a-uuid").Return(nil, notFoundError)
		var responseData map[string]any
		requestBuilder := testContext.readRequestBuilder("not-a-uuid").
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, "El elemento no existe", responseData["error"])
	})

	t.Run("reject an unauthenticated request", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		id := uuid.NewString()
		requestBuilder := testContext.readRequestBuilder(id).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		testContext.chunkRepository.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	})
}
