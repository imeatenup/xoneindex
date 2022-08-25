package intercept

import (
	"context"
	"net/http"

	"github.com/xnuc/xoneindex/log"

	"github.com/google/uuid"
)

type Trace struct{}

func (i *Trace) PreHandle(_ http.ResponseWriter, r *http.Request) bool {
	*r = *r.WithContext(context.WithValue(r.Context(), log.Trace{}, uuid.New().String()))
	log.Debugf(r.Context(), "Trace.PreHandle uri{%+v}", r.RequestURI)
	return true
}

func (i *Trace) PostHandle(_ http.ResponseWriter, r *http.Request) bool {
	log.Debugf(r.Context(), "Trace.PostHandle")
	return true
}
