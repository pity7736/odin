package accounting_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/accountsrepositoryfactory"
	"raiseexception.dev/odin/src/app"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/testutils"
)

const accountPath = "/accounts"

func TestHtmxShould(t *testing.T) {

	t.Run("create account when everything is ok", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body)

		response := testutils.GetJsonResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return unauthorized when request is anonymous", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithAnonymousSession()

		response := testutils.GetJsonResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})

	t.Run("return unauthorized when request sent cookie but session doesn't exists", func(t *testing.T) {
		accountingFactory := accountingrepositoryfactory.New()
		accountsFactory := accountsrepositoryfactory.New()
		user := userbuilder.New().Build()
		session := sessionmodel.New(user.ID())
		odinApp := app.NewFiberApplication(accountingFactory, accountsFactory)
		body := fmt.Sprintf(
			"name=%s&initial_balance=%s",
			"test",
			"10000",
		)
		requestBuilder := builders.NewRequestBuilder(accountsFactory).
			WithPath(accountPath).
			WithContentType(fiber.MIMEApplicationForm).
			WithPayload(body).
			WithSession(session)

		response := testutils.GetJsonResponseFromRequestBuilder(odinApp, requestBuilder)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, fiber.MIMETextHTMLCharsetUTF8, response.Header.Get("content-type"))
	})
}
