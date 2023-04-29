package main

import (
	"net/http"
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/sentry-go"
	"github.com/mariocoski/address-service/internal/app"
	"github.com/mariocoski/address-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewSentryLogger(config *config.Config) logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(
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
	logger.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.LogLevel)

	if err != nil {
		// log info level by default
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
	hook, err := logrus_sentry.NewSentryHook(config.SentryUrl, levels)
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true

	if err == nil {
		logger.Hooks.Add(hook)
	}
	return *logger
}

func main() {
	config := config.NewConfig()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryUrl,
		AttachStacktrace: true,
		// Debug: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		logrus.Fatalf("cannot initialise sentry %s", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	logger := NewSentryLogger(config)

	app := app.NewApplication(app.Dependencies{
		Config: config,
		Logger: logger,
	})

	server := &http.Server{
		Addr:    ":" + config.ApiPort,
		Handler: app,
	}

	logger.Infof("Listening on http://localhost:%v", config.ApiPort)

	serverErr := server.ListenAndServe()

	if serverErr != nil {
		logger.Fatalf("cannot start server: %v", serverErr.Error())
	}
}
