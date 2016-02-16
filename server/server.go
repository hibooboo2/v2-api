package server

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/auth"
	"github.com/rancher/v2-api/idformatter"

	"github.com/rancher/v2-api/vendor/github.com/gorilla/mux"
)

type Server struct {
	DB                 *sqlx.DB
	driver, driverName string
}

func New(driver, driverName string) (*Server, error) {
	db, err := sqlx.Open(driver, driverName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Server{
		driver:     driver,
		driverName: driverName,
		DB:         db,
	}, err
}

func (s *Server) namedQuery(query string, args map[string]interface{}) (*sqlx.Rows, error) {
	rows, err := s.DB.NamedQuery(query, args)
	return rows, err
}

func (s *Server) handleError(rw http.ResponseWriter, r *http.Request, err error) {
	var apiError *client.ServerApiError
	if e, ok := err.(*client.ServerApiError); ok {
		apiError = e
	} else {
		apiError = &client.ServerApiError{
			Type:    "error",
			Status:  500,
			Code:    "ServerError",
			Message: err.Error(),
		}
	}

	data, err := json.Marshal(apiError)
	if err == nil {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(apiError.Status)
		rw.Write(data)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("Fail to marshall: %v", err)
	}
}

func (s *Server) HandlerFunc(schemas *client.Schemas, f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return api.ApiHandlerFunc(schemas, func(rw http.ResponseWriter, r *http.Request) {
		logrus.Debugf("Vars avaliable: %#v", mux.Vars(r))
		ctx := api.GetApiContext(r)
		if ctx != nil {
			ctx.IDFormatter = idformatter.NewFormatter()
		}
		if err := auth.APIAuthenticator(rw, r); err != nil {
			s.handleError(rw, r, err)
		} else {
			if err := f(rw, r); err != nil {
				s.handleError(rw, r, err)
			}
		}
	})
}

func (s *Server) writeResponse(err error, r *http.Request, data interface{}) error {
	if err != nil {
		return err
	}
	api.GetApiContext(r).Write(data)
	return nil
}

func (s *Server) deobfuscate(r *http.Request, typeName string, id string) string {
	return api.GetApiContext(r).IDFormatter.ParseID(id)
}

func (s *Server) obfuscate(r *http.Request, typeName string, id string) string {
	ctx := api.GetApiContext(r)
	return ctx.IDFormatter.FormatID(id, typeName, ctx.Schemas)
}

func (s *Server) getClient(r *http.Request) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url: "http://localhost:8080/v1/projects/1a5/schemas",
	})
}

func (s *Server) parseInputParameters(r *http.Request) InputData {
	data := InputData{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	return data
}
