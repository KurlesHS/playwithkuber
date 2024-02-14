package logservermux

import (
	"log"
	"net/http"
)

type LogServeMux struct {
	http.ServeMux
}

func (mux *LogServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("http request: %s %s\n", r.Method, r.URL)
	mux.ServeMux.ServeHTTP(w, r)
}

func New() *LogServeMux {
	return new(LogServeMux)
}
