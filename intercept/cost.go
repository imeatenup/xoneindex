package intercept

import (
	"context"
	"net/http"
	"time"

	"github.com/xnuc/xoneindex/log"
)

type Cost struct{}

func (i *Cost) PreHandle(_ http.ResponseWriter, r *http.Request) bool {
	*r = *r.WithContext(context.WithValue(r.Context(), Cost{}, time.Now()))
	log.Debugf(r.Context(), "Cost.PreHandle")
	return true
}

func (i *Cost) PostHandle(_ http.ResponseWriter, r *http.Request) bool {
	log.Debugf(r.Context(), "Cost.PostHandle cost{%+v}", time.Since((r.Context().Value(Cost{})).(time.Time)))
	return true
}
