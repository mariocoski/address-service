package get_addresses_use_case

import (
	addresses "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
)

type GetAddressesUseCase struct {
	addressesRepository address_repo.AddressesRepository
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *GetAddressesUseCase {
	return &GetAddressesUseCase{
		addressesRepository: addressesRepository,
	}
}

func (uc *GetAddressesUseCase) GetAddresses() ([]*addresses.Address, error) {
	return uc.addressesRepository.GetAll()
}
