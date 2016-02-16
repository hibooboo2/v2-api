package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/rancher/v2-api/auth"
	"github.com/rancher/v2-api/router"
	"github.com/rancher/v2-api/server"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	listen := ":8899"
	logrus.Info("Starting Rancher V2 API on ", listen)
	config := mysql.Config{
		User:      "cattle",
		Passwd:    "cattle",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "cattle",
		Collation: "utf8_general_ci",
		ParseTime: true,
	}
	server, err := server.New("mysql", config.FormatDSN())
	if err != nil {
		logrus.Fatal(err)
	}

	err = auth.InitAuth(server.DB)
	if err != nil {
		logrus.Fatalf("Error when parsing key: %#v", err)
	}

	r := router.New(server)
	http.ListenAndServe(listen, r)
}
