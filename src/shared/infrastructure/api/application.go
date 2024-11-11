package api

import (
	"net/http"
)

type Application interface {
	Start() error
	Test(request *http.Request) (*http.Response, error)
}
