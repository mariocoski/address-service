package addresses_domain

import "errors"

type Address struct {
	Id           string `json:"id"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AddressLine3 string `json:"address_line_3"`
	City         string `json:"city"`
	County       string `json:"county"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
}

type AddressInitializer struct {
	AddressLine1 string `json:"address_line_1" validate:"required"`
	AddressLine2 string `json:"address_line_2"`
	AddressLine3 string `json:"address_line_3"`
	City         string `json:"city" validate:"required"`
	County       string `json:"county"`
	State        string `json:"state"`
	Postcode     string `json:"postcode" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

type AddressPatch struct {
	AddressLine1 *string `json:"address_line_1,omitempty"`
	AddressLine2 *string `json:"address_line_2,omitempty"`
	AddressLine3 *string `json:"address_line_3,omitempty"`
	City         *string `json:"city,omitempty"`
	County       *string `json:"county,omitempty"`
	State        *string `json:"state,omitempty"`
	Postcode     *string `json:"postcode,omitempty"`
	Country      *string `json:"country,omitempty"`
}

var ErrAddressNotFound = errors.New("address not found")
