package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
	jwt "gopkg.in/square/go-jose.v1"
	"time"
)

func unmarshalToken(tokenMarshaled string, token *Token) error {
	tokenParsed, err := jwt.ParseEncrypted(tokenMarshaled)
	if err != nil {
		return err
	}
	data, err := tokenParsed.Decrypt(jwtKey)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, token)
}

func authWithToken(w http.ResponseWriter, r *http.Request) error {
	token, err := getJWTString(r)
	if err != nil {
		return err
	}
	policy, err := getPolicy(token, r)
	if err != nil {
		return err
	}
	context := api.GetApiContext(r)
	if context == nil {
		return fmt.Errorf("No context found: %#v", r.URL)
	}
	context.Policy = policy
	return nil
}

func getJWTString(r *http.Request) (string, error) {
	token := ""

	for _, gotCookie := range r.Cookies() {
		if gotCookie.Name == "token" || gotCookie.Name == "TOKEN" {
			token = gotCookie.Value
		}
	}

	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.EqualFold(strings.SplitN(authHeader, "\\s", 2)[0], "bearer") {
			token = strings.Split(authHeader, "\\s")[1]
		}
	}

	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		return "", fmt.Errorf("Token not found")
	}
	return token, nil
}

func getPolicy(tokenKey string, r *http.Request) (*api.Policy, error) {
	gotToken := authToken{}
	query := `SELECT *
			FROM auth_token
			WHERE auth_token.key = ?`
	err := sqlxConn.Get(&gotToken, query, tokenKey)
	if err != nil {
		return nil, err

	}

	if time.Now().After(gotToken.Expires) {
		return nil, &expiredToken
	}

	logrus.Debugf("GotToken: %v expires: %v", gotToken.AccountID, gotToken.Expires)
	token := Token{}
	err = unmarshalToken(gotToken.Value, &token)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Token Expires: %v", time.Unix(token.Expires, 0))
	if time.Now().After(time.Unix(token.Expires, 0)) {
		return nil, &expiredToken
	}

	return createPolicy(token.RancherAccountID, token.Username, token.IdentitiesAsIds, r)
}

var expiredToken = client.ServerApiError{
	Type:    "error",
	Status:  401,
	Code:    "TokenExpired",
	Message: "Auth token is expired",
	Detail:  "Please relogin",
}
