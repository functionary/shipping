package ups

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"geo"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"
)

/*
PACKAGE TYPES:
If no container is specified, RAVE assumes UPS Package, i.e., type 02

Valid values:
00 = UNKNOWN;
01 = UPS Letter;
02 = Package;
03 = Tube;
04 = Pak;
21 = Express Box;
24 = 25KG Box;
25 = 10KG Box;
30 = Pallet;
2a = Small Express Box;
22b = Medium Express Box;
2c = Large Express Box
*/

type PackageType string

const (
	PackageUnknown          PackageType = "00"
	PackageUPSLetter        PackageType = "01"
	PackagePackage          PackageType = "02"
	PackageTube             PackageType = "03"
	PackagePak              PackageType = "04"
	PackageExpressBox       PackageType = "21"
	Package25KgBox          PackageType = "24"
	Package10KgBox          PackageType = "25"
	PackagePallet           PackageType = "30"
	PackageSmallExpressBox  PackageType = "2a"
	PackageMediumExpressBox PackageType = "2b"
	PackageLargeExpressBox  PackageType = "2c"
)

// Can this be made to 'const'?
var packageTypeNames map[PackageType]string = map[PackageType]string{
	PackageUnknown:          "Unknown Page",
	PackagePackage:          "Package",
	PackageUPSLetter:        "UPS Letter",
	PackageTube:             "Tube",
	PackagePak:              "Pak",
	PackageExpressBox:       "Express Box",
	Package25KgBox:          "25KG Box",
	Package10KgBox:          "10KG Box",
	PackagePallet:           "Pallet",
	PackageSmallExpressBox:  "Small Express Box",
	PackageMediumExpressBox: "Medium Express Box",
	PackageLargeExpressBox:  "Large Express Box"}

/* PICKUP TYPE:
Default value is 01.
Valid values are:
01 - Daily Pickup;
03 - Customer Counter;
06 - One Time Pickup;
07 - On Call Air;
19 - Letter Center;
20 - Air Service Center.
Refer to the Rate Chart table in Appendix C for rate type based on Pickup Type and Customer Classification Code.
*/

type PickupType string

const (
	PickupDaily            PickupType = "01"
	PickupCustomerCounter  PickupType = "03"
	PickupOneTime          PickupType = "06"
	PickupOnCallAir        PickupType = "07"
	PickupLetterCenter     PickupType = "19"
	PickupAirServiceCenter PickupType = "20"
)

/*
SERVICE:
Required for Rating and Ignored for Shopping.

Valid domestic values:
14 = Next Day Air Early AM,
01 = Next Day Air,
13 = Next Day Air Saver,
59 = 2nd Day Air AM,
02 = 2nd Day Air,
12 = 3 Day Select,
03 = Ground.

Valid international values:
11 = Standard,
07 = Worldwide Express,
54 = Worldwide Express Plus,
08 = Worldwide Expedited,
65 = Saver.

Valid Poland to Poland
Same Day values:
82 = UPS Today Standard,
83 = UPS Today Dedicated Courier,
84 = UPS Today Intercity,
85 = UPS Today Express,
86 = UPS Today Express Saver
*/

type Service string

const (
	// Domestic Services
	ServiceUSNextDayAirAM    Service = "14"
	ServiceUSNextDayAir      Service = "01"
	ServiceUSNextDayAirSaver Service = "13"
	ServiceUS2ndDayAirAM     Service = "59"
	ServiceUS2ndDayAir       Service = "02"
	ServiceUS3DaySelect      Service = "12"
	ServiceUSGround          Service = "03"

	// Internation Services
	ServiceIntlStandard         Service = "11"
	ServiceWorldwideExpress     Service = "07"
	ServiceWorldwideExpressPlus Service = "54"
	ServiceWorldwideExpedited   Service = "08"
	ServiceIntlSaver            Service = "65"
)

// Can this be made to 'const'?
var serviceNames map[Service]string = map[Service]string{
	ServiceUSNextDayAirAM:       "Next Day Air Early AM",
	ServiceUSNextDayAir:         "Next Day Air",
	ServiceUSNextDayAirSaver:    "Next Day Air Saver",
	ServiceUS2ndDayAirAM:        "2nd Day Air AM",
	ServiceUS2ndDayAir:          "2nd Day Air",
	ServiceUS3DaySelect:         "3 Day Select",
	ServiceUSGround:             "Ground",
	ServiceIntlStandard:         "International Standard",
	ServiceWorldwideExpress:     "Worldwide Express",
	ServiceWorldwideExpressPlus: "Worldwide Express Plus",
	ServiceWorldwideExpedited:   "Worldwide Expedited",
	ServiceIntlSaver:            "International Saver"}

type Credentials struct {
	AccessLicenseNumber string
	UserId              string
	Password            string
}

type Package struct {
	PackagingType struct {
		Code        PackageType
		Description string
	}
	Description     string // Merchandise description of package. 
	ReferenceNumber struct {
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
	LargePackageIndicator bool

	// Additional Handling:
	// The presence indicates additional handling is required.
	// The absence indicates no additional handling is required.
	AdditionalHandling bool

	Dimensions struct {
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

	IncludeDimensions      bool
	IncludeReferenceNumber bool
}

type Shipment struct {
	Shipper struct {
		ShipperNumber string
		Address       geo.Address
	}
	ShipTo   geo.Address
	ShipFrom geo.Address
	Service  struct {
		Code        Service
		Description string
	}
	PaymentInformation struct {
		Prepaid struct {
			BillShipper struct {
				AccountNumber string
			}
		}
	}
	Packages []Package
}

type Error struct {
	Source string
	Text   string
	Fatal  bool
}

func (e *Error) Error() string {
	var output string
	output = e.Source + ": " + e.Text
	if e.Fatal {
		output += " (fatal)"
	}
	return output
}

func (p *Package) validate() error {
	// Are dimensions required?
	// Required if Packaging Type is not
	// Letter, Express Tube, or Express Box;
	// Required for 'GB to GB' and 'Poland to Poland' shipments
	// Precision: 6.2
	// Assume inches (fixed in template).

	p.IncludeDimensions = false
	if (p.PackagingType.Code == PackageUPSLetter) || (p.PackagingType.Code == PackageTube) || (p.PackagingType.Code == PackageExpressBox) {
		p.IncludeDimensions = true
		if p.Dimensions.Width <= 0 {
			return errors.New("Package cannot have 0 width.")
		}
		if p.Dimensions.Height <= 0 {
			return errors.New("Package cannot have 0 height.")
		}
		if p.Dimensions.Length <= 0 {
			return errors.New("Package cannot have 0 length.")
		}
	}

	if p.PackageWeight.Weight > 150 {
		return errors.New("Package is too large (over 150lbs).")
	}
	if p.PackageWeight.Weight < 0.1 {
		return errors.New("Package is too small (under 0.1lbs).")
	}

	return nil
}


// Writes the standard UPS Access Request to the writer.
func accessRequest(w io.Writer, creds Credentials) error {
	// Parse and execute the templates.
	tdir := "./templates/xml/ups/"

	templates := []string{
		tdir + "accessrequest.xml",
	}
	page, err := template.ParseFiles(templates...)
	if err != nil {
		return errors.New("ups.accessRequest: Unable to parse templates:\n" + err.Error())
	}

	err = page.ExecuteTemplate(w, "accessrequest.xml", creds)
	if err != nil {
		return errors.New("UPS Test: Access request template execution failed:\n" + err.Error())
	}

	return nil
}

func send(data []byte, target string) ([]byte, error) {

	//	fmt.Printf("\n\n%s\n\n", data)

	bufin := bytes.NewBuffer(data)

	ioutil.WriteFile("request.xml", data, 0644)

	client := new(http.Client)

	request, err := http.NewRequest("POST", target, bufin)
	if err != nil {
		fmt.Println("ups.send: Error while creating XML request:\n", err.Error())
	}

	request.Header.Set("SOAPAction", "anAction")
	request.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	request.Header.Set("Content-Length", fmt.Sprintf("%d", bufin.Len()))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("ups.send: Error while sending XML request:\n", err.Error())
	}

	bufout := new(bytes.Buffer)
	r := bufio.NewReader(response.Body)
	defer response.Body.Close()

	line, _, err := r.ReadLine()
	for err == nil {
		bufout.Write(line)
		line, _, err = r.ReadLine()
	}
	if err != io.EOF {
		fmt.Println("ups.send: Error while reading response:\n", err.Error())
	}

	rawxml := bufout.Bytes()
	ioutil.WriteFile("response.xml", rawxml, 0644)

	return rawxml, nil

}

func NewPackage() *Package {
	p := new(Package)

	p.AdditionalHandling = false
	p.Dimensions.UnitOfMeasurement.Code = "IN"
	p.Dimensions.UnitOfMeasurement.Description = "Inches"
	p.IncludeDimensions = false
	p.IncludeReferenceNumber = false
	p.LargePackageIndicator = false
	p.PackageWeight.Weight = 0
	p.PackageWeight.UnitOfMeasurement.Code = "LBS"
	p.PackageWeight.UnitOfMeasurement.Description = "Pounds"
	p.PackagingType.Code = PackagePackage
	p.PackagingType.Description = packageTypeNames[PackagePackage]

	return p
}

func NewShipment() *Shipment {
	s := new(Shipment)

	s.Service.Code = ServiceUS2ndDayAir
	s.Service.Description = serviceNames[ServiceUSGround]

	return s
}
