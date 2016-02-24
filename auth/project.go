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

const memberQuery = `(external_id= ? and external_id_type= ? and state='active' and removed is null)`

func hasAccessToProject(projectID, usingAccount int64, isAdmin bool, identityIds []string) bool {
	logrus.Debugf("Project ID: %v Using %v ISADMIN %v Identities: %#v", projectID, usingAccount, isAdmin, len(identityIds))
	if len(identityIds) == 0 {
		return false
	}
	if usingAccount == projectID || isAdmin {
		logrus.Debugf("Is admin or self.")
		return true
	}
	query := `SELECT * from project_member where project_id= ? and ( `
	args := []interface{}{projectID}
	for _, identityID := range identityIds {
		spiltID := strings.Split(identityID, ":")
		args = append(args, spiltID[1])
		args = append(args, spiltID[0])
		query = query + memberQuery + " or "
	}
	query = strings.TrimSuffix(query, " or ") + " )"
	rows, err := sqlxConn.Query(query, args...)
	if err != nil {
		logrus.Debugf("Error getting members for project: %#v", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

type NoSpecifiedProject struct{}

func (e *NoSpecifiedProject) Error() string {
	return "No Project specified."
}

var NOProjectSpecified = &NoSpecifiedProject{}

func getProjectID(r *http.Request, formatter api.IDFormatter) (int64, error) {

	projectID := getProjectIDString(r)

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

func getProjectIDString(r *http.Request) string {
	projectID := ""

	theVars := mux.Vars(r)
	projectID, ok := theVars["envID"]
	if !ok {
		logrus.Debugf("Env from url: %#v PATH: %v All Vars: %#v", theVars["envID"], r.URL.Path, theVars)
	}

	if projectID == "" && strings.HasPrefix(r.URL.Path, "/v2/environments/") {
		projectID = strings.SplitN(strings.TrimPrefix(r.URL.Path, "/v2/environments/"), "/", 2)[0]
		logrus.Debugf("Using manual url parsing for project: %v", projectID)
	}

	if projectID == "" {
		logrus.Debugf("Env from header: %v", r.Header.Get(projectIDHeader))
		projectID = r.Header.Get(projectIDHeader)
	}

	if projectID == "" {
		logrus.Debugf("Env from query: %v", r.URL.Query().Get("projectId"))
		projectID = r.URL.Query().Get("projectId")
	}

	return projectID
}

const projectIDHeader = "X-API-Project-Id"
