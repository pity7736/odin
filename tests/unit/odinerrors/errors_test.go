package odinerrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/testutils"
)

var currentFilePath = fmt.Sprintf("%s/errors_test.go", testutils.GetTestPath())

func TestNewError(t *testing.T) {

	t.Run("it's built with default configuration", func(t *testing.T) {
		errorMessage := "some domain error"
		odinError := odinerrors.NewErrorBuilder(errorMessage).Build()
		var err *odinerrors.Error
		ok := errors.As(odinError, &err)

		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("%s:%s:19", errorMessage, currentFilePath), err.Error())
		assert.Equal(t, "", err.ExternalError())
	})

	t.Run("when has an external error message", func(t *testing.T) {
		errorMessage := "some domain error"
		externalMessage := "some external message"
		odinError := odinerrors.NewErrorBuilder(errorMessage).
			WithExternalMessage(externalMessage).
			Build()
		var err *odinerrors.Error
		ok := errors.As(odinError, &err)

		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("%s:%s:33", errorMessage, currentFilePath), err.Error())
		assert.Equal(t, externalMessage, err.ExternalError())
	})

	t.Run("when has a tag", func(t *testing.T) {
		odinError := odinerrors.NewErrorBuilder("some domain error").
			WithTag(odinerrors.DOMAIN).
			Build()
		var err *odinerrors.Error
		ok := errors.As(odinError, &err)

		assert.True(t, ok)
		assert.Equal(t, odinerrors.DOMAIN, err.Tag())
	})

	t.Run("when new error wraps base error", func(t *testing.T) {
		domainMessage := "some domain error"
		applicationMessage := "some application error"
		baseError := odinerrors.NewErrorBuilder(domainMessage).WithTag(odinerrors.DOMAIN).Build()
		externalError := "unexpected error"
		wrapError := odinerrors.NewErrorBuilder(applicationMessage).
			WithWrapped(baseError).
			WithExternalMessage(externalError).
			Build()
		var err *odinerrors.Error
		ok := errors.As(wrapError, &err)

		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("%s: %s:%s:56", applicationMessage, domainMessage, currentFilePath), err.Error())
		assert.Equal(t, externalError, err.ExternalError())
	})

	t.Run("when new error wraps base error with external message", func(t *testing.T) {
		domainMessage := "some domain error"
		domainExternalMessage := "account name cannot be empty"
		baseError := odinerrors.NewErrorBuilder(domainMessage).
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage(domainExternalMessage).
			Build()
		applicationExternalMessage := "unexpected error"
		applicationMessage := "some application error"
		errorWrapper := odinerrors.NewErrorBuilder(applicationMessage).
			WithWrapped(baseError).
			WithExternalMessage(applicationExternalMessage).
			Build()
		var err *odinerrors.Error
		ok := errors.As(errorWrapper, &err)

		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("%s: %s:%s:76", applicationMessage, domainMessage, currentFilePath), err.Error())
		assert.Equal(t, fmt.Sprintf("%s: %s", applicationExternalMessage, domainExternalMessage), err.ExternalError())
	})

	t.Run("when wrapped error is go error", func(t *testing.T) {
		errMessage := "some error"
		wrapperMessage := "some domain error"
		errorWrapper := odinerrors.NewErrorBuilder(wrapperMessage).
			WithWrapped(errors.New(errMessage)).
			Build()
		var err *odinerrors.Error
		ok := errors.As(errorWrapper, &err)

		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("%s: %s", wrapperMessage, errMessage), err.Error())
	})
}
