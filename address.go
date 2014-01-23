package shipping

import (
	"github.com/jackmanlabs/bucket"
)

type Address struct {
	NameFirst     string
	NameLast      string
	Company       string
	Street1       string
	Street2       string
	Street3       string
	City          string
	StateProvince string
	PostalCode    string
	Country       string
}

type AddressViewModel struct {
	NameFirst      string
	NameLast       string
	Company        string
	Street1        string
	Street2        string
	Street3        string
	City           string
	StateProvince  string
	StateProvinces []bucket.HTMLOption
	PostalCode     string
	Country        string
	Countries      []bucket.HTMLOption
}

func (addr *Address) ViewModel() *AddressViewModel {

	avm := new(AddressViewModel)
	avm.StateProvinces = make([]bucket.HTMLOption, 0)
	avm.Countries = make([]bucket.HTMLOption, 0)

	for _, state := range states {
		if state.Value == addr.StateProvince {
			state.Selected = true
		}
		avm.StateProvinces = append(avm.StateProvinces, state)
	}

	for _, country := range countries {
		if country.Value == addr.Country {
			country.Selected = true
		}
		avm.Countries = append(avm.Countries, country)
	}

	return avm
}
