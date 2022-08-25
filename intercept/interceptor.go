package intercept

import (
	"net/http"
)

type Intercept interface {
	PreHandle(http.ResponseWriter, *http.Request) bool
	PostHandle(http.ResponseWriter, *http.Request) bool
}

// Interceptor Intercept
// type Interceptor struct{}
