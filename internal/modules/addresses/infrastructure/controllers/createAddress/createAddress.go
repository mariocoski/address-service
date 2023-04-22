package create_address_controller

import (
	"encoding/json"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	useCase "github.com/mariocoski/address-service/internal/modules/addresses/application/createAddress"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/sirupsen/logrus"
)

type CreateAddressController struct {
	useCase useCase.CreateAddressUseCase
}

func NewController(useCase useCase.CreateAddressUseCase) *CreateAddressController {
	return &CreateAddressController{
		useCase: useCase,
	}
}

func (c *CreateAddressController) Handle(w http.ResponseWriter, r *http.Request) {

	var address domain.AddressInitializer

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&address)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(address)

	if err != nil {
		logrus.Error("CreateAddressController validation error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAddress, err := c.useCase.CreateAddress(address)

	// handle not found error with 429
	// https://github.com/jackc/pgx/issues/474#issuecomment-549397821

	if err != nil {
		logrus.Error("CreateAddressController useCase: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(createdAddress, "", "  ")
	// TODO: uncomment when done with development
	// response, err := json.Marshal(createdAddress, "", "  ")
	if err != nil {
		logrus.Error("CreateAddressController json marshalling: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
