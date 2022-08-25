package handle

import (
	"net/http"
)

type Handle interface {
	Handle(http.ResponseWriter, *http.Request) error
}
