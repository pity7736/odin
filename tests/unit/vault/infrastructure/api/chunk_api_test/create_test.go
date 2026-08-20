package chunk_api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountsinmemory "raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/apptest"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type chunkHandlerTestContext struct {
	application       app.Application
	userRepository    *accountsinmemory.InMemoryUserRepository
	sessionRepository *accountsinmemory.InMemorySessionRepository
	chunkRepository   *mocks.MockChunkRepository
}

func newChunkHandlerTestContext(t *testing.T) chunkHandlerTestContext {
	userRepository := accountsinmemory.NewInMemoryUserRepository()
	sessionRepository := accountsinmemory.NewInMemorySessionRepository()
	chunkRepository := mocks.NewMockChunkRepository(t)
	application := apptest.New().
		WithUserRepository(userRepository).
		WithSessionRepository(sessionRepository).
		WithChunkRepository(chunkRepository).
		Build()
	return chunkHandlerTestContext{
		application:       application,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		chunkRepository:   chunkRepository,
	}
}

func (self chunkHandlerTestContext) requestBuilder() *builders.RequestBuilder {
	return builders.NewRequestBuilder(self.userRepository, self.sessionRepository).
		WithPath("/api/v1/chunks").
		WithContentType(fiber.MIMEApplicationJSON)
}

func TestCreateChunkRestShould(t *testing.T) {
	t.Run("store a chunk for the authenticated owner", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		testContext.chunkRepository.EXPECT().GetByID(mock.Anything, owner.ID(), id).Return(nil, nil)
		var storedChunk *chunkmodel.EncryptedChunk
		testContext.chunkRepository.EXPECT().Add(mock.Anything, mock.Anything).Run(func(ctx context.Context, chunk *chunkmodel.EncryptedChunk) {
			storedChunk = chunk
		}).Return(nil)
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(fmt.Sprintf(`{"id": "%s", "content": "encrypted-content"}`, id)).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, id, responseData["id"])
		assert.NotNil(t, storedChunk)
		assert.Equal(t, id, storedChunk.ID())
		assert.Equal(t, owner.ID(), storedChunk.OwnerID())
		assert.Equal(t, "encrypted-content", storedChunk.Content())
		testContext.chunkRepository.AssertCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject a chunk whose id already exists for the owner", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		existingChunk, _ := chunkmodel.New(id, owner.ID(), "existing-content")
		testContext.chunkRepository.EXPECT().GetByID(mock.Anything, owner.ID(), id).Return(existingChunk, nil)
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(fmt.Sprintf(`{"id": "%s", "content": "encrypted-content"}`, id)).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusConflict, response.StatusCode)
		assert.Equal(t, "El elemento ya existe", responseData["error"])
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject a chunk with an empty id", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(`{"id": "", "content": "encrypted-content"}`).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Datos de solicitud inválidos", responseData["error"])
		testContext.chunkRepository.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject a chunk with empty content", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		id := uuid.NewString()
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(fmt.Sprintf(`{"id": "%s", "content": ""}`, id)).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Datos de solicitud inválidos", responseData["error"])
		testContext.chunkRepository.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject a chunk whose id is not a valid uuid", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		testContext.chunkRepository.EXPECT().GetByID(mock.Anything, owner.ID(), "not-a-uuid").Return(nil, nil)
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(`{"id": "not-a-uuid", "content": "encrypted-content"}`).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Datos de solicitud inválidos", responseData["error"])
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject a malformed body", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		owner := userbuilder.New().WithEmail("owner@example.com").Create(testContext.userRepository)
		var responseData map[string]any
		requestBuilder := testContext.requestBuilder().
			WithPayload(`{"id": "some-id" "content": "encrypted-content"}`).
			WithResponseData(&responseData).
			WithUser(owner)
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Datos de solicitud inválidos", responseData["error"])
		testContext.chunkRepository.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("reject an unauthenticated request", func(t *testing.T) {
		testContext := newChunkHandlerTestContext(t)
		id := uuid.NewString()
		requestBuilder := testContext.requestBuilder().
			WithPayload(fmt.Sprintf(`{"id": "%s", "content": "encrypted-content"}`, id)).
			WithAnonymousSession()
		response := testutils.GetJSONResponseFromRequestBuilder(testContext.application, requestBuilder)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		testContext.chunkRepository.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
		testContext.chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})
}
