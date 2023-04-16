package get_address_by_id_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	useCase "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddressById"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

type GetAddressByIdController struct {
	logger  logger.Logger
	useCase useCase.GetAddressByIdUseCase
}

func NewController(logger logger.Logger, useCase useCase.GetAddressByIdUseCase) *GetAddressByIdController {
	return &GetAddressByIdController{
		logger:  logger,
		useCase: useCase,
	}
}

func (c *GetAddressByIdController) Handle(w http.ResponseWriter, r *http.Request) {

	addressIdParam := chi.URLParam(r, "addressID")

	c.logger.Info(fmt.Sprintf(`GetAddressByIdController: received "addressId" url param: %v`, addressIdParam))

	if addressIdParam == "" {
		c.logger.Error("invalid addressId url param")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	address, err := c.useCase.GetById(addressIdParam)

	// handle not found error with 404
	// https://github.com/jackc/pgx/issues/474#issuecomment-549397821

	if err != nil {
		c.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(address, "", "  ")
	// TODO: uncomment when done with development
	// response, err := json.Marshal(addresss, "", "  ")
	if err != nil {
		c.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
