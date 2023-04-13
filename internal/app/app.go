package app

import (
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mariocoski/address-service/internal/config"
	get_addresses_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddresses"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	get_addresses_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/getAddresses"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

const API_PATH = "/api"
const ADDRESSES_PATH = "/addresses"

type Dependencies struct {
	Logger *logger.Logger
	Config *config.Config
}

func NewApplication(dependencies Dependencies) *chi.Mux {

	addressesRepository, err := address_repo.NewAddressesRepository(*dependencies.Config)

	getAddressesUseCase := get_addresses_use_case.NewUseCase(addressesRepository)
	getAddressesHandler := get_addresses_controller.NewGetAddressesController(*dependencies.Logger, getAddressesUseCase)

	if err != nil {
		log.Fatal("Cannot instantiate repo", err)
	}

	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route(fmt.Sprintf("%v%v", API_PATH, ADDRESSES_PATH), func(r chi.Router) {
		// GET /api/addresses
		r.Get("/", getAddressesHandler.Handle)
	})

	return mux
}
