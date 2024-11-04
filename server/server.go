package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {})

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
