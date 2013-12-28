package usps

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
	// "strconv"
	// "os"
)

/*
CONTAINER:
VARIABLE
FLAT RATE ENVELOPE
PADDED FLAT RATE ENVELOPE
LEGAL FLAT RATE ENVELOPE
SM FLAT RATE ENVELOPE
WINDOW FLAT RATE ENVELOPE
GIFT CARD FLAT RATE ENVELOPE
FLAT RATE BOX
SM FLAT RATE BOX
MD FLAT RATE BOX
LG FLAT RATE BOX
REGIONALRATEBOXA
REGIONALRATEBOXB
REGIONALRATEBOXC
RECTANGULAR
NONRECTANGULAR

Note: RECTANGULAR or NONRECTANGULAR must
be indicated when <Size>LARGE</Size>.
*/

type Container string

const (
	ContainerVariable            Container = "VARIABLE"
	ContainerFlatRateEnv         Container = "FLAT RATE ENVELOPE"
	ContainerFlatRateEnvPadded   Container = "PADDED FLAT RATE ENVELOPE"
	ContainerFlatRateEnvLegal    Container = "LEGAL FLAT RATE ENVELOPE"
	ContainerFlatRateEnvSmall    Container = "SM FLAT RATE ENVELOPE"
	ContainerFlatRateEnvWindow   Container = "WINDOW FLAT RATE ENVELOPE"
	ContainerFlatRateEnvGiftCard Container = "GIFT CARD FLAT RATE ENVELOPE"
	ContainerBoxFlatRate         Container = "FLAT RATE BOX"
	ContainerBoxFlatRateSmall    Container = "SM FLAT RATE BOX"
	ContainerBoxFlatRateMedium   Container = "MD FLAT RATE BOX"
	ContainerBoxFlatRateLarge    Container = "LG FLAT RATE BOX"
	ContainerRegionalRateBoxA    Container = "REGIONALRATEBOXA"
	ContainerRegionalRateBoxB    Container = "REGIONALRATEBOXB"
	ContainerRegionalRateBoxC    Container = "REGIONALRATEBOXC"
	ContainerRectangular         Container = "RECTANGULAR"
	ContainerNonrectangular      Container = "NONRECTANGULAR"
)

/*
FIRST CLASS MAIL TYPES:
LETTER
FLAT
PARCEL
POSTCARD
PACKAGE SERVICE
*/

type FirstClassType string

const (
	FirstClassLetter         FirstClassType = "LETTER"
	FirstClassFlat           FirstClassType = "FLAT"
	FirstClassParcel         FirstClassType = "PARCEL"
	FirstClassPostcard       FirstClassType = "POSTCARD"
	FirstClassPackageService FirstClassType = "PACKAGE SERVICE"
)

/*
SERVICE:
FIRST CLASS
FIRST CLASS COMMERCIAL
FIRST CLASS HFP COMMERCIAL
PRIORITY
PRIORITY COMMERCIAL
PRIORITY HFP COMMERCIAL
EXPRESS
EXPRESS COMMERCIAL
EXPRESS SH
EXPRESS SH COMMERCIAL
EXPRESS HFP
EXPRESS HFP COMMERCIAL
PARCEL
MEDIA
LIBRARY
ALL
ONLINE
*/

type Service string

const (
	ServiceFirstClass        Service = "FIRST CLASS"
	ServiceFirstClassComm    Service = "FIRST CLASS COMMERCIAL"
	ServiceFirstClassCommHFP Service = "FIRST CLASS HFP COMMERCIAL"
	ServicePriority          Service = "PRIORITY"
	ServicePriorityComm      Service = "PRIORITY COMMERCIAL"
	ServicePriorityCommHFP   Service = "PRIORITY HFP COMMERCIAL"
	ServiceExpress           Service = "EXPRESS"
	ServiceExpressComm       Service = "EXPRESS COMMERCIAL"
	ServiceExpressSH         Service = "EXPRESS SH"
	ServiceExpressCommSH     Service = "EXPRESS SH COMMERCIAL"
	ServiceExpressHFP        Service = "EXPRESS HFP"
	ServiceExpressCommHFP    Service = "EXPRESS HFP COMMERCIAL"
	ServiceParcel            Service = "PARCEL"
	ServiceMedia             Service = "MEDIA"
	ServiceLibrary           Service = "LIBRARY"
	ServiceAll               Service = "ALL"
	ServiceOnline            Service = "ONLINE"
)

type RatingServiceSelectionResponse struct {
	Response struct {
		ResponseStatusCode int
	}
	RatedShipment []RatedShipment
}

type RatedShipment struct {
	Service struct {
		Code int
	}
	TotalCharges struct {
		CurrencyCode  string
		MonetaryValue float64
	}
}

type Shipper struct {
	Address       geo.Address
	ShipperNumber string
}

type RateRequest struct {
	UserId   string
	Packages []Package
}

type Estimate struct {
	Description string // Verbal description of the estimate.
	Service     Service
	Cost        float64
}

type Package struct {
	Service        Service
	FirstClassType FirstClassType
	Container      Container
	ZipTo          string
	ZipFrom        string

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

	// SIZE:
	// REGULAR: Package dimensions are 12" or less;
	// LARGE: Any package dimension is larger than 12".
	Size string

	// Template Controls:
	IsLarge      bool
	IsFirstClass bool
}

func (p *Package) validate() error {
	p.Size = "REGULAR"
	p.IsLarge = false
	if (p.Width > 12) || (p.Height > 12) || (p.Length > 12) {
		p.Size = "LARGE"
		p.IsLarge = true
	}

	if (p.Service == ServiceFirstClass) || (p.Service == ServiceFirstClass) || (p.Service == ServiceFirstClass) {
		p.IsFirstClass = true
	}

	return nil
}

func Rate(request *RateRequest) (Estimate, error) {

	buf := new(bytes.Buffer)

	var estimate Estimate

	if err := requestRate(buf, request); err != nil {
		return estimate, errors.New("ups.Rate: Rate request failed:\n" + err.Error())
	}

	rawxml, err := send(buf.Bytes())
	if err != nil {
		return estimate, errors.New("ups.Shop: Data send failed:\n" + err.Error())
	}

	// var tree interface{}
	var response RatingServiceSelectionResponse

	err = xml.Unmarshal(rawxml, &response)
	if err != nil {
		return estimate, errors.New("ups.Shop: XML unmarshalling failed:\n" + err.Error())
	}

	var estimates []Estimate
	for _, value := range response.RatedShipment {
		var e Estimate
		e.Service = Service(value.Service.Code)
		e.Cost = value.TotalCharges.MonetaryValue
		estimates = append(estimates, e)
	}

	return estimate, nil
}

func Shop(request *RateRequest) ([]Estimate, error) {
	buf := new(bytes.Buffer)

	if err := requestRate(buf, request); err != nil {
		return nil, errors.New("ups.Shop: Rate request failed:\n" + err.Error())
	}

	rawxml, err := send(buf.Bytes())
	if err != nil {
		return nil, errors.New("ups.Shop: Data send failed:\n" + err.Error())
	}

	// var tree interface{}
	var response RatingServiceSelectionResponse

	err = xml.Unmarshal(rawxml, &response)
	if err != nil {
		return nil, errors.New("ups.Shop: XML unmarshalling failed:\n" + err.Error())
	}

	var estimates []Estimate
	for _, value := range response.RatedShipment {
		var e Estimate
		e.Service = Service(value.Service.Code)
		e.Cost = value.TotalCharges.MonetaryValue
		estimates = append(estimates, e)
	}

	return estimates, nil
}

// This does the actual processing of the UPS Rate Request. Rate() and Shop() are both front-ends to this function.
func requestRate(w io.Writer, req *RateRequest) error {
	// Parse and execute the templates.
	tdir := "./templates/usps/"

	templates := []string{
		tdir + "raterequest.xml"}
	page, err := template.ParseFiles(templates...)
	if err != nil {
		return errors.New("USPS Test: Unable to parse templates:\n" + err.Error())
	}

	err = page.ExecuteTemplate(w, "raterequest.xml", req)
	if err != nil {
		return errors.New("USPS Test: Rate request template execution failed:\n" + err.Error())
	}

	return nil
}

func send(data []byte) ([]byte, error) {
	fmt.Printf("\n\n%s\n\n", data)

	//	target := "https://secure.shippingapis.com/ShippingAPITest.dll"
	target := "http://testing.shippingapis.com//ShippingAPITest.dll"

	ioutil.WriteFile("request.xml", data, 0644)

	client := new(http.Client)

	values := url.Values{}
	values.Add("API", "RateV2")
	values.Add("XML", fmt.Sprintf("%s", data))

	response, err := client.PostForm(target, values)
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

	fmt.Printf("\n\n%s\n\n", rawxml)

	return rawxml, nil

}
