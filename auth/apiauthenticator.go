package auth

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
)

func APIAuthenticator(w http.ResponseWriter, r *http.Request) error {
	ctx := api.GetApiContext(r)
	if ctx == nil {
		return errors.New("No context found for apirequest.")
	}

	if err := authenticate(w, r); err != nil {
		if e, ok := err.(*client.ServerApiError); ok {
			return e
		}
		logrus.Debugf("Failed to auth with go. Attempting cattle: %#v", err)
		if err = forwardAuthToCattle(w, r); err != nil {
			if e, ok := err.(*client.ServerApiError); ok {
				return e
			}
			logrus.Debugf("Failed to auth by proxing to cattle: %#v", err)
			return &client.ServerApiError{
				Type:    "error",
				Status:  http.StatusUnauthorized,
				Code:    "Unauthorized",
				Message: "Unauthorized",
			}
		}
	}
	logrus.Debugf("Policy is: %v %v %v %v", ctx.Policy.AccountID, ctx.Policy.AuthenticatedAsAccountID,
		ctx.Policy.Username, len(ctx.Policy.Identities))
	return nil
}

func authenticate(w http.ResponseWriter, r *http.Request) error {

	if err := authWithToken(w, r); err != nil {
		if e, ok := err.(*client.ServerApiError); ok {
			return e
		}
		logrus.Debugf("Failed auth with token: %#v", err)
		return basicAuth(w, r)
	}
	return nil
}
