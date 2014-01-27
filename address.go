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

	// This literal declaration is much prettier than explicit assignment to the
	// struct fields.

	avm := &AddressViewModel{
		NameFirst:      addr.NameFirst,
		NameLast:       addr.NameLast,
		Company:        addr.Company,
		Street1:        addr.Street1,
		Street2:        addr.Street2,
		Street3:        addr.Street3,
		City:           addr.City,
		StateProvince:  addr.StateProvince,
		StateProvinces: make([]bucket.HTMLOption, 0),
		PostalCode:     addr.PostalCode,
		Country:        addr.Country,
		Countries:      make([]bucket.HTMLOption, 0),
		Phone:          addr.Phone,
		Fax:            addr.Fax,
		Email:          addr.Email,
	}

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
