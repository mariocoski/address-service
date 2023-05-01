package addresses_domain

import (
	"errors"

	"github.com/google/uuid"
)

type Address struct {
	Id           uuid.UUID `json:"id" bson:"_id"`
	AddressLine1 string    `json:"address_line_1" bson:"address_line_1"`
	AddressLine2 string    `json:"address_line_2" bson:"address_line_2"`
	AddressLine3 string    `json:"address_line_3" bson:"address_line_3"`
	City         string    `json:"city" bson:"city"`
	County       string    `json:"county" bson:"county"`
	State        string    `json:"state" bson:"state"`
	Postcode     string    `json:"postcode" bson:"postcode"`
	Country      string    `json:"country" bson:"country"`
	CreatedAt    string    `json:"created_at" bson:"created_at"`
	UpdatedAt    string    `json:"updated_at" bson:"updated_at"`
}

type AddressInitializer struct {
	AddressLine1 string `json:"address_line_1" bson:"address_line_1"  validate:"required"`
	AddressLine2 string `json:"address_line_2" bson:"address_line_2"`
	AddressLine3 string `json:"address_line_3" bson:"address_line_3"`
	City         string `json:"city" bson:"city"  validate:"required"`
	County       string `json:"county" bson:"county"`
	State        string `json:"state" bson:"state"`
	Postcode     string `json:"postcode" bson:"postcode"  validate:"required"`
	Country      string `json:"country" bson:"country"  validate:"required"`
}

type AddressDocument struct {
	Id           string `json:"_id" bson:"_id"  validate:"required"`
	AddressLine1 string `json:"address_line_1" bson:"address_line_1"  validate:"required"`
	AddressLine2 string `json:"address_line_2" bson:"address_line_2"`
	AddressLine3 string `json:"address_line_3" bson:"address_line_3"`
	City         string `json:"city" bson:"city"  validate:"required"`
	County       string `json:"county" bson:"county"`
	State        string `json:"state" bson:"state"`
	Postcode     string `json:"postcode" bson:"postcode"  validate:"required"`
	Country      string `json:"country" bson:"country"  validate:"required"`
	CreatedAt    string `json:"created_at" bson:"created_at"`
	Deleted      string `json:"deleted_at" bson:"deleted_at"`
	UpdatedAt    string `json:"updated_at" bson:"updated_at"`
}

type AddressPatch struct {
	AddressLine1 *string `json:"address_line_1, bson:"address_line_1, omitempty"`
	AddressLine2 *string `json:"address_line_2, bson:"address_line_2, omitempty"`
	AddressLine3 *string `json:"address_line_3, bson:"address_line_3, omitempty"`
	City         *string `json:"city, bson:"city, omitempty"`
	County       *string `json:"county, bson:"county, omitempty"`
	State        *string `json:"state, bson:"state, omitempty"`
	Postcode     *string `json:"postcode, bson:"postcode, omitempty"`
	Country      *string `json:"country, bson:"country, omitempty"`
}

var ErrAddressNotFound = errors.New("address not found")
