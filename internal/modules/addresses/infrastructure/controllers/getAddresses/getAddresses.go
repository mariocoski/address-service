package get_addresses_controller

import (
	"encoding/json"
	"net/http"

	get_addresses_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddresses"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

type GetAddressesController struct {
	logger  logger.Logger
	useCase *get_addresses_use_case.GetAddressesUseCase
}

func NewGetAddressesController(logger logger.Logger, useCase *get_addresses_use_case.GetAddressesUseCase) *GetAddressesController {
	return &GetAddressesController{
		logger:  logger,
		useCase: useCase,
	}
}

func (c *GetAddressesController) Handle(w http.ResponseWriter, r *http.Request) {
	addresses, err := c.useCase.GetAddresses()
	if err != nil {
		c.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(addresses)
	if err != nil {
		c.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
