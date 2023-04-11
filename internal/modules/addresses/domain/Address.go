package addresses_domain

type Address struct {
	Id           string `json:"id"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	City         string `json:"city"`
	County       string `json:"county,omitempty"`
	State        string `json:"state,omitempty"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
}
