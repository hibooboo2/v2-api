package auth

import (
	"database/sql"
	"time"

	"github.com/rancher/go-rancher/api"
)

type Token struct {
	AccountID        string   `json:"account_id"`
	RancherAccountID int64    `json:"rancher_account_id"`
	AccountType      string   `json:"external_account_type"`
	IdentitiesAsIds  []string `json:"idList"`
	Username         string   `json:"username"`
	Expires          int64    `json:"exp"`
}

type authToken struct {
	ID        int64     `db:"id"`
	AccountID int64     `db:"account_id"`
	Created   time.Time `db:"created"`
	Expires   time.Time `db:"expires"`
	Key       string    `db:"key"`
	Value     string    `db:"value"`
	Version   string    `db:"version"`
	Provider  string    `db:"provider"`
}

type Data struct {
	ID      int64        `db:"id"`
	Name    string       `db:"name"`
	Visible bool         `db:"visible"`
	Value   sql.RawBytes `db:"value"`
}

type Policies struct {
	Data []api.Policy `json:"data"`
}
