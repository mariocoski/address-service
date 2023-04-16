package get_address_by_id_use_case

import (
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	address_repo "github.com/mariocoski/address-service/internal/modules/addresses/domain/repositories"
)

type GetAddressByIdUseCase interface {
	GetById(addressId string) (domain.Address, error)
}

type GetAddressById struct {
	addressesRepository address_repo.AddressesRepository
}

func NewUseCase(addressesRepository address_repo.AddressesRepository) *GetAddressById {
	return &GetAddressById{
		addressesRepository: addressesRepository,
	}
}

func (uc *GetAddressById) GetById(addressId string) (domain.Address, error) {
	return uc.addressesRepository.GetById(addressId)
}
