package main

import (
	"net/http"

	"github.com/rancher/v2-api/router"
	"github.com/rancher/v2-api/server"

	"github.com/Sirupsen/logrus"
)

func main() {
	listen := ":8899"
	logrus.Info("Starting Rancher V2 API on ", listen)
	server := server.New()
	r := router.New(server)
	http.ListenAndServe(listen, r)
}
