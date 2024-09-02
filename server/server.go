package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, `{"status": "OK"}`)
	})

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
