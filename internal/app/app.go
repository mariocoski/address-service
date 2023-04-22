package app

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mariocoski/address-service/internal/config"
	create_address_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/createAddress"
	delete_address_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/deleteAddress"
	get_address_by_id_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddressById"
	get_addresses_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddresses"
	update_address_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/updateAddress"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	create_address_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/createAddress"
	delete_address_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/deleteAddress"
	get_address_by_id_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/getAddressById"
	get_addresses_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/getAddresses"
	update_address_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/updateAddress"
	"github.com/mariocoski/address-service/internal/shared/http/middlewares"
	"github.com/sirupsen/logrus"
)

const API_PATH = "/api"
const ADDRESSES_PATH = "/addresses"

type Dependencies struct {
	Config *config.Config
}

func NewApplication(dependencies Dependencies) *chi.Mux {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dependencies.Config.SentryUrl,
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

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	addressesRepository, err := address_repo.NewAddressesRepository(*dependencies.Config)

	getAddressesUseCase := get_addresses_use_case.NewUseCase(addressesRepository)
	getAddressesController := get_addresses_controller.NewController(getAddressesUseCase)

	getAddressByIdUseCase := get_address_by_id_use_case.NewUseCase(addressesRepository)
	getAddressByIdController := get_address_by_id_controller.NewController(getAddressByIdUseCase)

	createAddressUseCase := create_address_use_case.NewUseCase(addressesRepository)
	createAddressController := create_address_controller.NewController(createAddressUseCase)

	deleteAddressUseCase := delete_address_use_case.NewUseCase(addressesRepository)
	deleteAddressController := delete_address_controller.NewController(deleteAddressUseCase)

	updateAddressUseCase := update_address_use_case.NewUseCase(addressesRepository)
	updateAddressController := update_address_controller.NewController(updateAddressUseCase)

	if err != nil {
		logrus.Fatal("cannot instantiate repo", err)
	}

	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// Important: Chi has a middleware stack and thus it is important to put the
	// Sentry handler on the appropriate place. If using middleware.Recoverer,
	// the Sentry middleware must come afterwards (and configure it with
	// Repanic: true).
	mux.Use(sentryMiddleware.Handle)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middlewares.Authenticate(*dependencies.Config))

	mux.Route(fmt.Sprintf("%v%v", API_PATH, ADDRESSES_PATH), func(r chi.Router) {
		r.Post("/", createAddressController.Handle)
		r.Get("/{addressID}", getAddressByIdController.Handle)
		r.Get("/", getAddressesController.Handle)
		r.Patch("/{addressID}", updateAddressController.Handle)
		r.Delete("/{addressID}", deleteAddressController.Handle)
	})

	return mux
}
