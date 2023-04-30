package get_addresses_controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"

	"github.com/sirupsen/logrus"
)

type GetAddressesController struct {
	addressesRepository address_repo.AddressesRepository
}

func NewController(addressesRepository address_repo.AddressesRepository) *GetAddressesController {
	return &GetAddressesController{
		addressesRepository,
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

	addressesPaginated, err := c.addressesRepository.GetAllPaginated(currentPage, perPage)

	if err != nil {
		logrus.Error(err.Error())
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(addressesPaginated)

	if err != nil {
		logrus.Error(err.Error())
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
