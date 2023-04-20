package delete_address_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	useCase "github.com/mariocoski/address-service/internal/modules/addresses/application/deleteAddress"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

type DeleteAddressController struct {
	logger  logger.Logger
	useCase useCase.DeleteAddressUseCase
}

func NewController(logger logger.Logger, useCase useCase.DeleteAddressUseCase) *DeleteAddressController {
	return &DeleteAddressController{
		logger:  logger,
		useCase: useCase,
	}
}

func (c *DeleteAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	c.logger.Info(fmt.Sprintf(`DeleteAddressController: received "addressId" url param: %v`, addressIdParam))

	if addressIdParam == "" {
		c.logger.Error("invalid addressId url param")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	addressId, err := c.useCase.DeleteAddress(addressIdParam)
	log.Println("addressId", addressId)

	// https://github.com/jackc/pgx/issues/474#issuecomment-549397821
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		c.logger.Error(fmt.Sprintf("cannot delete address by id: %v, error: %v", addressId, err))
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
		c.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
