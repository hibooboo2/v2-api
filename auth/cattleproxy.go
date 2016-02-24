package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
)

func forwardAuthToCattle(w http.ResponseWriter, r *http.Request) error {
	var url string
	projectID := getProjectIDString(r)
	if projectID != "" {
		url = fmt.Sprintf("http://localhost:8080/v1/projects/%v/apirequestpolicy?%v", projectID, r.URL.RawQuery)
	} else {
		url = "http://localhost:8080/v1/apirequestpolicy?" + r.URL.RawQuery
	}
	logrus.Debugf("Using : %v", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", r.Header.Get("Cookie"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return &client.ServerApiError{
			Type:    "error",
			Status:  resp.StatusCode,
			Code:    "Not sure",
			Message: "Need to get from resp body",
			Detail:  "Still need to get",
		}
	}

	var policy api.Policy
	var policyList Policies
	err = json.Unmarshal(body, &policyList)
	if err != nil {
		return err
	}

	policy = policyList.Data[0]
	logrus.Debugf("Policy user is: %#v account %v", policy.Username, policy.AccountID)

	if len(policy.Identities) > 0 {
		api.GetApiContext(r).Policy = &policy
		return nil
	}

	return fmt.Errorf("Implement cattle forward for authentication.")
}
