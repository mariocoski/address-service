package delete_address_controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	useCase "github.com/mariocoski/address-service/internal/modules/addresses/application/deleteAddress"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/sirupsen/logrus"
)

type DeleteAddressController struct {
	useCase useCase.DeleteAddressUseCase
}

func NewController(useCase useCase.DeleteAddressUseCase) *DeleteAddressController {
	return &DeleteAddressController{
		useCase: useCase,
	}
}

func (c *DeleteAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	logrus.Infof(`DeleteAddressController: received "addressId" url param: %v`, addressIdParam)

	if addressIdParam == "" {
		logrus.Error("invalid addressId url param")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	addressId, err := c.useCase.DeleteAddress(addressIdParam)

	// https://github.com/jackc/pgx/issues/474#issuecomment-549397821
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logrus.Errorf("cannot delete address by id: %v, error: %v", addressId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := struct {
		Id string `json:"id"`
	}{
		Id: addressId,
	}

	responseJson, err := json.MarshalIndent(response, "", "  ")
	// TODO: uncomment when done with development
	// responseJson, err := json.Marshal(response, "", "  ")
	if err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
