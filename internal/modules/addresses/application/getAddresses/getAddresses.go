package get_addresses_use_case

import (
	addresses "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
	"github.com/mariocoski/address-service/internal/shared/core/pagination"
	"github.com/mariocoski/address-service/internal/shared/logger"
)

type GetAddressesUseCase struct {
	addressesRepository address_repo.AddressesRepository
	logger              logger.Logger
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *GetAddressesUseCase {
	return &GetAddressesUseCase{
		addressesRepository: addressesRepository,
	}
}

type Request struct {
	CurrentPage int
	PerPage     int
}

func (uc *GetAddressesUseCase) GetAddresses(request Request) (pagination.PaginationResult[addresses.Address], error) {
	return uc.addressesRepository.GetAllPaginated(request.CurrentPage, request.PerPage)
}
