package create_address_controller

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	validator "github.com/go-playground/validator/v10"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	"github.com/sirupsen/logrus"
)

type CreateAddressController struct {
	addressesRepository address_repo.AddressesRepository
}

func NewController(addressesRepository address_repo.AddressesRepository) *CreateAddressController {
	return &CreateAddressController{
		addressesRepository,
	}
}

func (c *CreateAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	var address domain.AddressInitializer

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&address)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sentry.CaptureException(err)
		return
	}

	err = validator.New().Struct(address)

	if err != nil {
		logrus.Error("CreateAddressController validation error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAddress, err := c.addressesRepository.Save(address)

	// handle not found error with 429
	// https://github.com/jackc/pgx/issues/474#issuecomment-549397821

	if err != nil {
		logrus.Error("CreateAddressController repo: " + err.Error())
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(createdAddress)

	if err != nil {
		logrus.Error("CreateAddressController json marshalling: " + err.Error())
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
