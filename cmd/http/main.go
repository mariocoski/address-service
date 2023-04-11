package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mariocoski/address-service/internal/app"
	"github.com/mariocoski/address-service/internal/config"
	repositories "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

const HOST = "localhost:7000"

func main() {
	config := config.NewConfig()
	logger := logger.NewLogger()

	repoDependencies := repositories.AddresssRepoDependencies{
		RepoType: "postgres",
		Config:   config,
	}
	addressesRepository, err := repositories.NewAddressesRepository(&repoDependencies)

	if err != nil {
		log.Fatal("Cannot instantiate repo", err)
	}

	dependencies := app.Dependencies{
		Logger:              &logger,
		Config:              config,
		AddressesRepository: addressesRepository,
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
