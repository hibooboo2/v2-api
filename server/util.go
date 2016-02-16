package server

import (
	"net/http"

	"github.com/rancher/go-rancher/api"
)

func (s *Server) getAccountID(r *http.Request) int64 {
	ctx := api.GetApiContext(r)
	if ctx == nil {
		return 0
	}
	policy := ctx.Policy
	if policy == nil {
		return 0
	}
	return policy.AccountID
}
