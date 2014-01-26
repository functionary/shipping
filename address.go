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
	Phone         string
	Fax           string
	Email         string
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
	Phone          string
	Fax            string
	Email          string
}

func (addr *Address) ViewModel() *AddressViewModel {

	avm := new(AddressViewModel)
	avm.StateProvinces = make([]bucket.HTMLOption, 0)
	avm.Countries = make([]bucket.HTMLOption, 0)

	avm.NameFirst = addr.NameFirst
	avm.NameLast = addr.NameLast
	avm.Company = addr.Company
	avm.Street1 = addr.Street1
	avm.Street2 = addr.Street2
	avm.Street3 = addr.Street3
	avm.City = addr.City
	avm.StateProvince = addr.StateProvince
	avm.PostalCode = addr.PostalCode
	avm.Country = addr.Country
	avm.Phone = addr.Phone
	avm.Fax = addr.Fax
	avm.Email = addr.Email

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
