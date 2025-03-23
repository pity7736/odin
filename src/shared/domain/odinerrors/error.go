package odinerrors

import (
	"fmt"
	"strings"
)

const sep = ": "

type Error struct {
	tag      Tag
	message  string
	location string
	external string
	err      *Error
}

func New(message string) *Error {
	return &Error{message: message}
}

func (self *Error) Error() string {
	current := self
	builder := strings.Builder{}
	builder.Grow(len(current.message))
	for current != nil {
		if current.location != "" {
			builder.WriteString(fmt.Sprintf("%s:%s", current.message, current.location))
		} else {
			builder.WriteString(current.message)
		}
		current = current.err
		if current != nil {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}

func (self *Error) ExternalError() string {
	current := self
	builder := strings.Builder{}
	builder.Grow(len(current.external))
	for current != nil {
		builder.WriteString(current.external)
		current = current.err
		if current != nil && current.external != "" {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}

func (self *Error) Tag() Tag {
	return self.tag
}
