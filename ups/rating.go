package ups

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"text/template"
)

type RatingServiceSelectionRequest struct {
	Credentials Credentials

	Request struct {
		TransactionReference struct {
			CustomerContext string
		}
		RequestAction string
		RequestOption string
	}
	PickupType struct {
		Code        PickupType
		Description string
	}
	Shipment Shipment

	// Template Controls:
	IsShopping bool
}

type RatingServiceSelectionResponse struct {
	Response struct {
		ResponseStatusCode int
	}
	RatedShipment []struct {
		Service struct {
			Code int
		}
		TotalCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}
	}
}

type Estimate struct {
	Description string // Verbal description of the estimate.
	Service     Service
	Cost        float64
}

// This does the actual processing of the UPS Rate Request. Rate() and Shop() are both front-ends to this function.
func requestRate(w io.Writer, req *RatingServiceSelectionRequest) error {
	// Parse and execute the templates.
	tdir := "./templates/xml/ups/"

	templates := []string{
		tdir + "raterequest.xml"}
	page, err := template.ParseFiles(templates...)
	if err != nil {
		return errors.New("UPS Test: Unable to parse templates:\n" + err.Error())
	}

	err = page.ExecuteTemplate(w, "raterequest.xml", req)
	if err != nil {
		return errors.New("UPS Test: Rate request template execution failed:\n" + err.Error())
	}

	return nil
}

func Shop(request *RatingServiceSelectionRequest) ([]Estimate, error) {
	buf := new(bytes.Buffer)

	request.Request.RequestOption = "Shop"
	request.IsShopping = true

	if err := accessRequest(buf, request.Credentials); err != nil {
		return nil, errors.New("ups.Shop: Access request failed:\n" + err.Error())
	}

	if err := requestRate(buf, request); err != nil {
		return nil, errors.New("ups.Shop: Rate request failed:\n" + err.Error())
	}

	target := "https://wwwcie.ups.com/ups.app/xml/Rate"
	rawxml, err := send(buf.Bytes(), target)
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
		e.Description = serviceNames[e.Service]
		estimates = append(estimates, e)
	}

	return estimates, nil
}

func Rate(request *RatingServiceSelectionRequest) (Estimate, error) {

	buf := new(bytes.Buffer)

	request.Request.RequestOption = "Shop"
	request.IsShopping = false

	var estimate Estimate

	if err := accessRequest(buf, request.Credentials); err != nil {
		return estimate, errors.New("ups.Rate: Access request failed:\n" + err.Error())
	}

	if err := requestRate(buf, request); err != nil {
		return estimate, errors.New("ups.Rate: Rate request failed:\n" + err.Error())
	}

	target := "https://wwwcie.ups.com/ups.app/xml/Rate"
	rawxml, err := send(buf.Bytes(), target)
	if err != nil {
		return estimate, errors.New("ups.Shop: Data send failed:\n" + err.Error())
	}

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
		e.Description = serviceNames[e.Service]
		estimates = append(estimates, e)
	}

	if len(estimates) == 0 {
		return estimate, errors.New("ups.Shop: No results!:\n" + err.Error())
	}

	return estimates[0], nil
}
