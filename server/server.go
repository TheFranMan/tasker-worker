package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	r := mux.NewRouter()
	r.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {})
	r.Handle("/metrics", promhttp.Handler())

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
