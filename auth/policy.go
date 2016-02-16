package auth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
)

func createPolicy(accountID int64, username string, identitiesAsIds []string, r *http.Request) (*api.Policy, error) {
	formatter := api.GetApiContext(r).IDFormatter

	projectID, err := getProjectID(r, formatter)

	if err != nil && err != NOProjectSpecified {
		return nil, err
	}
	logrus.Debugf("Project picked: %v", projectID)
	if err == NOProjectSpecified {
		return &api.Policy{
			AuthenticatedAsAccountID: accountID,
			AccountID:                accountID,
			Identities:               []api.Identity{},
			Username:                 username,
		}, nil
	}
	canAccessProject := hasAccessToProject(projectID, accountID, isAdmin(accountID), identitiesAsIds)
	if !canAccessProject {
		return nil, &client.ServerApiError{
			Status: 403,
			Code:   "Forbidden",
			Type:   "error",
		}
	}
	return &api.Policy{
		AuthenticatedAsAccountID: accountID,
		AccountID:                projectID,
		Identities:               []api.Identity{},
		Username:                 username,
	}, nil
}
