package auth

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var sqlxConn *sqlx.DB

var jwtKey interface{}

func InitAuth(db *sqlx.DB) error {
	sqlxConn = db
	keyData := sqlxConn.QueryRow("SELECT value FROM data where name='host.api.key'")

	data := []byte{}
	err := keyData.Scan(&data)

	if err != nil {
		return err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("No block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	jwtKey = privateKey
	return nil
}
