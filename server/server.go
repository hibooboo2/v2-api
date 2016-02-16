package server

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rancher/go-rancher/api"

	"github.com/rancher/go-rancher/client"
)

type Server struct {
	db *sqlx.DB
}

func New() *Server {
	return &Server{}
}

func (s *Server) handleError(rw http.ResponseWriter, r *http.Request, err error) {
}

func (s *Server) HandlerFunc(schemas *client.Schemas, f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return api.ApiHandlerFunc(schemas, func(rw http.ResponseWriter, r *http.Request) {
		if err := f(rw, r); err != nil {
			s.handleError(rw, r, err)
		}
	})
}
