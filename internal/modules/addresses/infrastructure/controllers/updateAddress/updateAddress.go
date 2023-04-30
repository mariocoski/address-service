package update_address_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	"github.com/sirupsen/logrus"
)

type UpdateAddressController struct {
	addressesRepository address_repo.AddressesRepository
}

func NewController(addressesRepository address_repo.AddressesRepository) *UpdateAddressController {
	return &UpdateAddressController{
		addressesRepository,
	}
}

func (c *UpdateAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	logrus.Infof(`UpdateAddressController: received "addressId" url param: %v`, addressIdParam)

	if addressIdParam == "" {
		logrus.Error("UpdateAddressController: missing addressId url param")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var addressPatch domain.AddressPatch

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&addressPatch)

	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hasUpdates := false
	if addressPatch.AddressLine1 != nil || addressPatch.AddressLine2 != nil || addressPatch.AddressLine3 != nil ||
		addressPatch.City != nil || addressPatch.County != nil || addressPatch.State != nil || addressPatch.Postcode != nil ||
		addressPatch.Country != nil {
		hasUpdates = true
	}

	if !hasUpdates {
		http.Error(w, "You must provide at least one property to update address", http.StatusBadRequest)
		return
	}

	if addressPatch.AddressLine1 != nil && *addressPatch.AddressLine1 == "" {
		http.Error(w, "address_line_1 cannot be blank when explicilty defined", http.StatusBadRequest)
		return
	}

	if addressPatch.City != nil && *addressPatch.City == "" {
		http.Error(w, "city cannot be blank when explicilty defined", http.StatusBadRequest)
		return
	}

	if addressPatch.Postcode != nil && *addressPatch.Postcode == "" {
		http.Error(w, "postcode cannot be blank when explicilty defined", http.StatusBadRequest)
		return
	}

	if addressPatch.Country != nil && *addressPatch.Country == "" {
		http.Error(w, "country cannot be blank when explicilty defined", http.StatusBadRequest)
		return
	}

	updatedAddres, err := c.addressesRepository.Update(addressIdParam, addressPatch)

	if err != nil {
		sentry.CaptureException(err)
		if errors.Is(err, domain.ErrAddressNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logrus.Error(fmt.Sprintf("cannot get address by id: %v, error: %v", addressIdParam, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedAddres)
	if err != nil {
		sentry.CaptureException(err)
		logrus.Error("UpdateAddressController: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
