package http

import (
	"fmt"
	"net/http"
)

func NewAdminHandler() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("", makeHandler(nodeAdminHandler))
	return s
}

func makeHandler(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func nodeAdminHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
