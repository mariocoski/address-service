package update_address_use_case

import (
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
)

type UpdateAddressUseCase interface {
	UpdateAddress(addressId string, address domain.AddressPatch) (domain.Address, error)
}

type UpdateAddress struct {
	addressesRepository address_repo.AddressesRepository
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *UpdateAddress {
	return &UpdateAddress{
		addressesRepository: addressesRepository,
	}
}

func (uc *UpdateAddress) UpdateAddress(addressId string, patch domain.AddressPatch) (domain.Address, error) {
	return uc.addressesRepository.Update(addressId, patch)
}
