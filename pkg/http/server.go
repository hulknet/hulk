package http

import (
	"fmt"
	"net/http"

	"github.com/kotfalya/hulk/pkg/utils"
)

type Server struct {
	public  *http.Server
	private *http.Server
	admin   *http.Server
}

func NewServer(address string) *Server {
	h := http.NewServeMux()
	h.HandleFunc("/", handler)

	public := &http.Server{
		Handler:   h,
		Addr:      address,
		TLSConfig: utils.GenerateTLSConfig(),
	}

	h2 := http.NewServeMux()
	h2.HandleFunc("/", handler)
	h2.Handle("/node", http.StripPrefix("/node", h))
	private := &http.Server{
		Handler:   h2,
		Addr:      address,
		TLSConfig: utils.GenerateTLSConfig(),
	}

	return &Server{public, private, private}
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
