package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/api"
)

const memberQuery = `(external_id='%v' and external_id_type='%v' and state='active' and removed is null)`

func hasAccessToProject(projectID, usingAccount int64, isAdmin bool, identityIds []string) bool {
	logrus.Debugf("Project ID: %v Using %v ISADMIN %v Identities: %#v", projectID, usingAccount, isAdmin, len(identityIds))
	if len(identityIds) == 0 {
		return false
	}
	if usingAccount == projectID || isAdmin {
		logrus.Debugf("Is admin or self.")
		return true
	}
	query := fmt.Sprintf(`SELECT * from project_member where project_id='%v' and ( `, projectID)
	for _, identityID := range identityIds {
		spiltID := strings.Split(identityID, ":")
		query = query + fmt.Sprintf(memberQuery, spiltID[1], spiltID[0]) + " or "
	}
	query = strings.TrimSuffix(query, " or ") + " )"
	logrus.Debugf("Query is: %s", query)
	rows, err := sqlxConn.Query(query)
	if err != nil {
		return false
	}
	return rows.Next()
}

type NoSpecifiedProject struct{}

func (e *NoSpecifiedProject) Error() string {
	return "No Project specified."
}

var NOProjectSpecified = &NoSpecifiedProject{}

func getProjectID(r *http.Request, formatter api.IDFormatter) (int64, error) {
	projectID := ""

	theVars := mux.Vars(r)
	logrus.Debugf("Env from url: %v PATH: %v All Vars: %#v", theVars["envID"], r.URL.Path, theVars)
	projectID = theVars["envID"]

	if projectID == "" && strings.HasPrefix(r.URL.Path, "/v2/environments/") {
		projectID = strings.SplitN(strings.TrimPrefix(r.URL.Path, "/v2/environments/"), "/", 2)[0]
	}

	if projectID == "" {
		logrus.Debugf("Env from header: %v", r.Header.Get(projectIDHeader))
		projectID = r.Header.Get(projectIDHeader)
	}

	if projectID == "" {
		logrus.Debugf("Env from query: %v", r.URL.Query().Get("projectId"))
		projectID = r.URL.Query().Get("projectId")
	}

	if projectID == "" {
		return 0, NOProjectSpecified
	}

	parsedProjectID := formatter.ParseID(projectID)

	if parsedProjectID == "" {
		return 0, fmt.Errorf("Invalid projectIdFormat. ProjectId Header format is incorrect.")
	}

	logrus.Debugf("ProjectID = %v", parsedProjectID)

	accountID, err := strconv.ParseInt(parsedProjectID, 10, 64)
	if err != nil {
		return 0, err
	}

	if accountExists(accountID) {
		return accountID, nil
	}

	return 0, errors.New("Project not found.")
}

const projectIDHeader = "X-API-Project-Id"
