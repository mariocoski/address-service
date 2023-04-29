package get_address_by_id_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"

	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/sirupsen/logrus"
)

type GetAddressByIdController struct {
	addressesRepository address_repo.AddressesRepository
}

func NewController(addressesRepository address_repo.AddressesRepository) *GetAddressByIdController {
	return &GetAddressByIdController{
		addressesRepository,
	}
}

func (c *GetAddressByIdController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	logrus.Info(fmt.Sprintf(`GetAddressByIdController: received "addressId" url param: %v`, addressIdParam))

	if addressIdParam == "" {
		logrus.Errorf("missing addressId url param")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	address, err := c.addressesRepository.GetById(addressIdParam)

	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logrus.Errorf("cannot get address by id: %v, error: %v", addressIdParam, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(address, "", "  ")
	// TODO: uncomment when done with development
	// response, err := json.Marshal(addresss, "", "  ")
	if err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
