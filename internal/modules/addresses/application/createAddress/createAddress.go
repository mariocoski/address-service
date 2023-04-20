package create_address_use_case

import (
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
)

type CreateAddressUseCase interface {
	CreateAddress(address domain.AddressInitializer) (domain.Address, error)
}

type CreateAddress struct {
	addressesRepository address_repo.AddressesRepository
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *CreateAddress {
	return &CreateAddress{
		addressesRepository: addressesRepository,
	}
}

func (uc *CreateAddress) CreateAddress(address domain.AddressInitializer) (domain.Address, error) {

	return uc.addressesRepository.Save(address)
}
