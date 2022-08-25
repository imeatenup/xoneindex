package handle

import (
	"net/http"

	"github.com/xnuc/xoneindex/intercept"
)

type Service interface {
	http.Handler
}

type Server struct {
	Interceptor []intercept.Intercept
	Handler     Handle
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cnt := 0
	defer func(w http.ResponseWriter, r *http.Request) {
		for idx := cnt - 1; idx >= 0; idx-- {
			if !s.Interceptor[idx].PostHandle(w, r) {
				return
			}
		}
	}(w, r)
	for idx := 0; idx < len(s.Interceptor); idx++ {
		if !s.Interceptor[idx].PreHandle(w, r) {
			return
		}
		cnt++
	}
	err := s.Handler.Handle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
