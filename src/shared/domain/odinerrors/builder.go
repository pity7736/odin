package odinerrors

import (
	"errors"
	"fmt"
	"runtime"
)

type OdinErrorBuilder struct {
	tag      Tag
	message  string
	location string
	external string
	err      *Error
}

func NewErrorBuilder(message string) *OdinErrorBuilder {
	return &OdinErrorBuilder{message: message}
}

func (self *OdinErrorBuilder) WithExternalMessage(message string) *OdinErrorBuilder {
	self.external = message
	return self
}

func (self *OdinErrorBuilder) WithTag(tag Tag) *OdinErrorBuilder {
	self.tag = tag
	return self
}

func (self *OdinErrorBuilder) WithWrapped(err error) *OdinErrorBuilder {
	var odinError *Error
	ok := errors.As(err, &odinError)
	if !ok {
		odinError = &Error{message: err.Error()}
	}
	self.err = odinError
	return self
}

func (self *OdinErrorBuilder) Build() error {
	err := New(self.message)
	err.external = self.external
	err.err = self.err
	if self.err == nil {
		err.location = getLocation()
	} else {
		if self.tag == UNKNOWN {
			self.tag = self.err.Tag()
		}
	}
	err.tag = self.tag
	return err
}

func getLocation() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d", frame.File, frame.Line)
}
