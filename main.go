package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/rancher/v2-api/router"
	"github.com/rancher/v2-api/server"
	"flag"
	"os"
	"fmt"
)


var (
	VERSION = "HEAD"
)

func main() {
	processCmdLineFlags()


	listen := ":8899"
	logrus.Info("Starting Rancher V2 API on ", listen)
	config := mysql.Config{
		User:      "cattle",
		Passwd:    "cattle",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "cattle",
		Collation: "utf8_general_ci",
	}
	server, err := server.New("mysql", config.FormatDSN())
	if err != nil {
		logrus.Fatal(err)
	}
	r := router.New(server)
	http.ListenAndServe(listen, r)
}

func processCmdLineFlags() {
	// Define command line flags
	logLevel := flag.String("loglevel", "info", "Set the default loglevel (default:info) [debug|info|warn|error]")
	version := flag.Bool("v", false, "read the version of the v2-api")
	output := flag.String("o", "", "set the output file to write logs into, default is stdout")

	flag.Parse()

	if *output != "" {
		var f *os.File
		if _, err := os.Stat(*output); os.IsNotExist(err) {
			f, err = os.Create(*output)
			if err != nil {
				fmt.Printf("could not create file=%s for logging, err=%v\n", *output, err)
				os.Exit(1)
			}
		} else {
			var err error
			f, err = os.OpenFile(*output, os.O_RDWR|os.O_APPEND, 0)
			if err != nil {
				fmt.Printf("could not open file=%s for writing, err=%v\n", *output, err)
				os.Exit(1)
			}
		}
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(f)
	}

	if *version {
		fmt.Printf("v2-api\t gitcommit=%s\n", VERSION)
		os.Exit(0)
	}

	// Process log level.  If an invalid level is passed in, we simply default to info.
	if parsedLogLevel, err := logrus.ParseLevel(*logLevel); err == nil {
		logrus.WithFields(logrus.Fields{
			"logLevel": *logLevel,
		}).Info("Setting log level")
		logrus.SetLevel(parsedLogLevel)
	}
}
