package ups

/*
In order to work with the UPS servers peaceably, there's some trickery in play.
UPS specifies a number of XML types that may be included zero or one times
based on the type of transaction, but the UPS servers are very picky. For
example, the Rate Request wants most of those structures to be included zero
times, or else an error will be thrown.

To deal with this problem, the request structs often are defined with slices of
sub-structs. Because of the way the XML processor works, if there are zero
objects in the slice, no XML is rendered for the sub-struct.

The alternative is to create distinct kinds of sub-structs depending on the kind
of request.

At some point, we might think about validating for the correct number of 
sub-structs before requests are actually made. 
*/

import ()

type AddressType struct {
	AddressLine1                string
	AddressLine2                string `xml:",omitempty"`
	AddressLine3                string `xml:",omitempty"`
	City                        string
	StateProvinceCode           string
	PostalCode                  string
	CountryCode                 string
	ResidentialAddressIndicator string `xml:",omitempty"`
}

type ShipperType struct {
	Name          string
	ShipperNumber string
	Address       AddressType
}

type ShipToType struct {
	CompanyName string
	Address     AddressType
}

type ShipFromType struct {
	CompanyName string
	Address     AddressType
}

type ShipmentType struct {
	Description string
	Shipper     ShipperType
	ShipTo      ShipToType
	ShipFrom    ShipFromType
	Service     struct {
		Code ServiceCode
	}
	DocumentsOnly      string
	NumOfPieces        string
	PaymentInformation struct {
		Prepaid struct {
			BillShipper struct {
				AccountNumber string
			}
		}
	}
	Packages               []PackageType `xml:"Package"`
	ShipmentServiceOptions []struct {
		OnCallAir struct {
			Schedule struct {
				PickupDay int
				Method    int
			}
		}
	}
	RateInformation []struct {
		NegotiatedRatesIndicator string `xml:",omitempty"`
		RateChartIndicator       string `xml:",omitempty"`
	}
	InvoiceLineTotal []struct {
		CurrencyCode  string
		MonetaryValue float64
	}
	ItemizedChargesRequestedIndicator string `xml:",omitempty"`
}

type PackageType struct {
	PackagingType struct {
		Code        PackagingTypeCode
		Description string
	}
	Description     string // Merchandise description of package.
	ReferenceNumber []struct {
		Code  string
		Value string
	}
	PackageWeight struct {
		// Weight:
		// Assume pounds (fixed in template).
		// Precision: 6.1
		// Valid Range: 0.1-150.0
		UnitOfMeasurement struct {
			Code        string
			Description string
		}
		Weight float64
	}
	LargePackageIndicator bool `xml:",omitempty"`

	// Additional Handling:
	// The presence indicates additional handling is required.
	// The absence indicates no additional handling is required.
	AdditionalHandling bool `xml:",omitempty"`

	Dimensions []struct {
		// Width/Length/Height:
		// Required if Packaging Type is not
		// Letter, Express Tube, or Express Box;
		// Required for 'GB to GB' and 'Poland to Poland' shipments
		// Precision: 6.2

		UnitOfMeasurement struct {
			Code        string
			Description string
		}
		Width  float64
		Height float64
		Length float64
	}
}

type UnitOfMeasurementType struct {
	Code        UnitOfMeasurementCode
	Description string
}
