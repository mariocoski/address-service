package get_addresses_controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	get_addresses_use_case "github.com/mariocoski/address-service/internal/modules/addresses/application/getAddresses"
	"github.com/sirupsen/logrus"
)

type GetAddressesController struct {
	useCase *get_addresses_use_case.GetAddressesUseCase
}

func NewController(useCase *get_addresses_use_case.GetAddressesUseCase) *GetAddressesController {
	return &GetAddressesController{
		useCase: useCase,
	}
}

const DEFAULT_NUMBER_OF_ADDRESSES_PER_PAGE = 10
const MAX_NUMBER_OF_ADDRESSES_PER_PAGE = 50

func (c *GetAddressesController) Handle(w http.ResponseWriter, r *http.Request) {
	currentPage := 1
	perPage := DEFAULT_NUMBER_OF_ADDRESSES_PER_PAGE

	pageParam := r.URL.Query().Get("page")
	perPagaParam := r.URL.Query().Get("per_page")

	logrus.Infof(`GetAddressesController: received "page" query param with value: %v`, pageParam)
	logrus.Infof(`GetAddressesController: received: "per_page" query param with value: %v`, perPagaParam)

	if pageParam != "" {
		pageParamAsInt, err := strconv.Atoi(pageParam)
		if err == nil && pageParamAsInt > 0 {
			currentPage = pageParamAsInt
		}
	}

	if perPagaParam != "" {
		perPagaParamAsInt, err := strconv.Atoi(perPagaParam)
		if err == nil && perPagaParamAsInt >= 1 && perPagaParamAsInt <= MAX_NUMBER_OF_ADDRESSES_PER_PAGE {
			perPage = perPagaParamAsInt
		}
	}

	request := get_addresses_use_case.Request{
		CurrentPage: currentPage,
		PerPage:     perPage,
	}

	addressesPaginated, err := c.useCase.GetAddresses(request)

	if err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(addressesPaginated, "", "  ")
	// TODO: uncomment when done with development
	// response, err := json.Marshal(addressesPaginated, "", "  ")
	if err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
