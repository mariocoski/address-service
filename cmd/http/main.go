package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mariocoski/address-service/internal/app"
	"github.com/mariocoski/address-service/internal/config"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

const HOST = "localhost:7000"

func main() {
	config := config.NewConfig()
	logger := logger.NewLogger()

	dependencies := app.Dependencies{
		Logger: &logger,
		Config: config,
	}

	app := app.NewApplication(dependencies)

	server := &http.Server{
		Addr:    HOST,
		Handler: app,
	}

	fmt.Printf("Listening on http://%v", HOST)

	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatal("cannot start server", serverErr)
	}
}
