package constants

import (
	"fmt"
	"strings"
)

type CategoryType int8

const (
	_                    = iota
	EXPENSE CategoryType = iota
	INCOME
)

func NewFromString(value string) (CategoryType, error) {
	value = strings.ToLower(value)
	switch value {
	case "expense":
		return EXPENSE, nil
	case "income":
		return INCOME, nil
	default:
		return 0, fmt.Errorf("%s is an invalid category type", value)
	}
}

func (self CategoryType) String() string {
	if self == EXPENSE {
		return "expense"
	}
	return "income"
}
