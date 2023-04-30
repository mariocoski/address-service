package delete_address_controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	address_repo "github.com/mariocoski/address-service/internal/domain/repositories"

	domain "github.com/mariocoski/address-service/internal/domain"
	"github.com/sirupsen/logrus"
)

type DeleteAddressController struct {
	addressesRepository address_repo.AddressesRepository
}

func NewController(addressesRepository address_repo.AddressesRepository) *DeleteAddressController {
	return &DeleteAddressController{
		addressesRepository,
	}
}

func (c *DeleteAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	logrus.Infof(`DeleteAddressController: received "addressId" url param: %v`, addressIdParam)

	if addressIdParam == "" {
		logrus.Error("DeleteAddressController: invalid addressId url param")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	addressId, err := c.addressesRepository.Delete(addressIdParam)

	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sentry.CaptureException(err)
		logrus.WithField("addressId", addressId).WithError(err).Error("DeleteAddressController: cannot delete address by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := struct {
		Id string `json:"id"`
	}{
		Id: addressId,
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		logrus.Error(err.Error())
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
