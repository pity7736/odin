package constants_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounting/domain/constants"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		name     string
		input    constants.CategoryType
		expected string
	}{
		{
			name:     "when value is expense",
			input:    constants.EXPENSE,
			expected: "expense",
		},
		{
			name:     "when value is income",
			input:    constants.INCOME,
			expected: "income",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, testCase.input.String())
	}
}

func TestNewFromStringSuccess(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected constants.CategoryType
	}{
		{
			name:     "when value is expense in lower case",
			input:    "expense",
			expected: constants.EXPENSE,
		},
		{
			name:     "when value is expense upper case",
			input:    "EXPENSE",
			expected: constants.EXPENSE,
		},
		{
			name:     "when value is income in lower case",
			input:    "income",
			expected: constants.INCOME,
		},
		{
			name:     "when value is income in upper case",
			input:    "INCOME",
			expected: constants.INCOME,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			categoryType, err := constants.NewFromString(testCase.input)

			assert.Nil(t, err)
			assert.Equal(t, testCase.expected.String(), categoryType.String())
		})
	}
}

func TestNewFromStringWhenValueIsInvalid(t *testing.T) {
	value := "invalid category type"
	categoryType, err := constants.NewFromString(value)

	assert.Equal(t, constants.CategoryType(0), categoryType)
	assert.Equal(t, fmt.Errorf("%s is an invalid category type", value), err)
}
