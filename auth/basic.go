package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
)

const rancherID = "rancher_id"

func basicAuth(w http.ResponseWriter, r *http.Request) error {
	username, password, authFound := r.BasicAuth()
	if !authFound {
		return fmt.Errorf("No basic auth found.")
	}
	accountID, ok := getAccountID(username, password)
	if !ok {
		return errors.New("No account found.")
	}
	policy, err := createPolicy(accountID, "", []string{fmt.Sprintf("%v:%v", rancherID, accountID)}, r)

	if err != nil {
		return err
	}

	api.GetApiContext(r).Policy = policy

	return nil
}

const accountQuery = `SELECT account.id, credential.secret_value, account.kind
						FROM credential
						JOIN account on account.id = credential.account_id
						WHERE credential.public_value = :public_value`
const shaPrefix = `SHA256:`

type accountAndPass struct {
	ID          int64  `db:"id"`
	SecretValue string `db:"secret_value"`
	Kind        string `db:"kind"`
}

func getAccountID(publicValue, secretValue string) (int64, bool) {
	logrus.Debugf("User: %#v, Pass: %#v", publicValue, secretValue)
	accountPass := accountAndPass{}
	logrus.Debugf("Query: %v %v", accountQuery, publicValue)

	rows, err := sqlxConn.NamedQuery(accountQuery, map[string]interface{}{
		"public_value": publicValue,
	})
	if err != nil {
		logrus.Errorf("Error from retrieving accountPass: %#v", err)
		return 0, false
	}
	logrus.Debugf("Rows are: %#v", rows)

	if rows.Next() {
		if err = rows.Scan(&accountPass.ID, &accountPass.SecretValue, &accountPass.Kind); err != nil {
			logrus.Errorf("Errored on scanning row: %#v", err)
			return 0, false
		}
	}

	if rows.Next() {
		logrus.Debugf("Multiple apikeys not allowed.")
		return 0, false
	}
	logrus.Debugf("Account pass: %#v", accountPass)

	if strings.HasPrefix(accountPass.SecretValue, shaPrefix) {
		split := strings.Split(accountPass.SecretValue[len(shaPrefix):], ":")
		hasher := sha256.New()
		hasher.Write([]byte(split[0]))
		hasher.Write([]byte(secretValue))
		calculatedHash := hex.EncodeToString(hasher.Sum(nil))
		if calculatedHash == split[1] {
			logrus.Debugf("Authed as %v successfully.", accountPass.ID)
			return accountPass.ID, true
		}
	}
	return 0, false
}
