package shipping

// This is a more generic form, used primary for output.
type Estimate struct {
	Name     string
	Provider Carrier
	Service  string
	Price    float64
}

// A generic form, primarily used for quick transfer of data into carrier-specific APIs.
type Package struct {
	// Weight:
	// Units are ounces.
	// Cannot be greater than 70lbs (1120oz).
	Weight float64

	// Width/Length/Height:
	// Units are inches.
	// Required when SIZE is LARGE.
	Width  float64
	Height float64
	Length float64
}

type Carrier string

const (
	FedEx Carrier = "FedEx"
	UPS   Carrier = "UPS"
	USPS  Carrier = "USPS"
)
