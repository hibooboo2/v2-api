package auth

import "github.com/Sirupsen/logrus"

const getAccountByID = `SELECT kind
						FROM account
						WHERE account.id = ?`

func isAdmin(accountID int64) bool {
	row := sqlxConn.QueryRow(getAccountByID, accountID)
	var kind string
	if err := row.Scan(&kind); err != nil {
		logrus.Errorf("Error checking is admin: %v", err)
		return false
	}
	logrus.Debugf("Account %v is kind %v", accountID, kind)
	return kind == "admin" || kind == "superadmin"
}

func accountExists(accountID int64) bool {
	row := sqlxConn.QueryRowx(getAccountByID, accountID)
	var kind string
	if err := row.Scan(&kind); err != nil {
		logrus.Debugf("Account %v doesn't exists", accountID)
		return false
	}
	logrus.Debugf("Account %v is exists %v", accountID, kind)
	return true
}
