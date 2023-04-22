package main

import (
	"net/http"
	"os"

	"github.com/mariocoski/address-service/internal/app"
	"github.com/mariocoski/address-service/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {

	config := config.NewConfig()

	logrus.SetFormatter(
		&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function_name",
				logrus.FieldKeyFile:  "path_name",
			},
		},
	)
	logrus.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.LogLevel)

	if err != nil {
		// log info level by default
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	app := app.NewApplication(app.Dependencies{
		Config: config,
	})

	server := &http.Server{
		Addr:    ":" + config.ApiPort,
		Handler: app,
	}

	logrus.Infof("Listening on http://localhost:%v", config.ApiPort)

	serverErr := server.ListenAndServe()

	if serverErr != nil {
		logrus.Fatalf("cannot start server: %v", serverErr.Error())
	}
}
