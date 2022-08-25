package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xnuc/xoneindex/handle"
	"github.com/xnuc/xoneindex/intercept"
	"github.com/xnuc/xoneindex/layout"
)

func main() {
	server := &http.Server{
		Addr:         "127.0.0.1:9000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	mux := http.NewServeMux()
	mux.Handle("/", &handle.Server{
		Interceptor: []intercept.Intercept{&intercept.Trace{}, &intercept.Cost{}, &intercept.OauthClient{}},
		Handler:     &layout.Index{}})
	server.Handler = mux
	log.Fatal(server.ListenAndServe())
}
