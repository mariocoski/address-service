package delete_address_use_case

import (
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
)

type DeleteAddressUseCase interface {
	DeleteAddress(addressId string) (string, error)
}

type DeleteAddress struct {
	addressesRepository address_repo.AddressesRepository
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *DeleteAddress {
	return &DeleteAddress{
		addressesRepository: addressesRepository,
	}
}

func (uc *DeleteAddress) DeleteAddress(addressId string) (string, error) {
	return uc.addressesRepository.Delete(addressId)
}
