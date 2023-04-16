package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mariocoski/address-service/internal/config"
	get_address_by_id_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddressById"
	get_addresses_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddresses"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	get_address_by_id_controller "github.com/mariocoski/address-service/internal/modules/addresses/infrastructure/controllers/getAddressById"
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
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dependencies.Config.SentryUrl,
		Debug: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	addressesRepository, err := address_repo.NewAddressesRepository(*dependencies.Config)

	getAddressesUseCase := get_addresses_use_case.NewUseCase(addressesRepository)
	getAddressesController := get_addresses_controller.NewController(*dependencies.Logger, getAddressesUseCase)

	getAddressByIdUseCase := get_address_by_id_use_case.NewUseCase(addressesRepository)
	getAddressByIdController := get_address_by_id_controller.NewController(*dependencies.Logger, getAddressByIdUseCase)

	if err != nil {
		log.Fatal("Cannot instantiate repo", err)
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

	mux.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		hub := sentry.GetHubFromContext(r.Context())
		hub.CaptureException(errors.New("test error"))
	})
	mux.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("server panic")
	})

	mux.Route(fmt.Sprintf("%v%v", API_PATH, ADDRESSES_PATH), func(r chi.Router) {
		// GET /api/addresses
		r.Get("/{addressID}", getAddressByIdController.Handle)
		r.Get("/", getAddressesController.Handle)

		// GET /api/addresses/{addressId}
	})

	return mux
}
