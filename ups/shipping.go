package ups

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"geo"
	"text/template"
)

type ShipmentConfirmRequest struct {
	Credentials Credentials
	ShipTo      geo.Address
	ShipFrom    geo.Address
	Shipper     geo.Address

	Request struct {
		TransactionReference struct {
			CustomerContext string
		}
		RequestAction string
		RequestOption string
	}

	LabelSpecification struct {
		LabelPrintMethod struct {
			Code        string
			Description string
		}
		HTTPUserAgent    string
		LabelImageFormat struct {
			Code        string
			Description string
		}
		LabelStockSize struct {
			// In inches, whole numbers only.
			Width  int
			Height int
		}
	}
	Shipment Shipment
	Service  struct {
		Code        string
		Description string
	}
	Packages []Package
}

// XML Structure
type ShipmentConfirmResponse struct {
	Response struct {
		ResponseStatusCode        int
		ResponseStatusDescription string
		TransactionReference      struct {
			CustomerContext string
		}
		Error []struct {
			ErrorSeverity    string
			ErrorCode        int
			ErrorDescription string
		}
	}
	ShipmentCharges struct {
		TransportationCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}
		ServiceOptionsCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}
		TotalCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}
	}
	NegotiatedRates struct {
		NetSummaryCharges struct {
			GrandTotal struct {
				CurrencyCode  string
				MonetaryValue float64
			}
		}
	}
	BillingWeight struct {
		UnitOfMeasurement struct {
			Code string
		}
		Weight float64
	}
	ShipmentIdentificationNumber string
	ShipmentDigest               string
}

type ShipmentAcceptRequest struct {
	Credentials Credentials

	Request struct {
		TransactionReference struct {
			CustomerContext string
		}
		RequestAction string
		RequestOption string
	}
	ShipmentDigest string
}

type ShipmentAcceptResponse struct {
	Response struct {
		ResponseStatusCode        int
		ResponseStatusDescription string
		TransactionReference      struct {
			CustomerContext string
		}
		Error []struct {
			ErrorSeverity    string
			ErrorCode        int
			ErrorDescription string
		}
	}
	ShipmentResults struct {
		ShipmentCharges struct {
			TransportationCharges struct {
				CurrencyCode  string
				MonetaryValue float64
			}
			ServiceOptionsCharges struct {
				CurrencyCode  string
				MonetaryValue float64
			}
			TotalCharges struct {
				CurrencyCode  string
				MonetaryValue float64
			}
		}
		NegotiatedRates struct {
			NetSummaryCharges struct {
				GrandTotal struct {
					CurrencyCode  string
					MonetaryValue float64
				}
			}
		}
		BillingWeight struct {
			UnitOfMeasurement struct {
				Code string
			}
			Weight float64
		}
		ShipmentIdentificationNumber string
		PickupRequestNumber          string
		PackageResults               []struct {
			TrackingNumber        string
			ServiceOptionsCharges struct {
				CurrencyCode  string
				MonetaryValue float64
			}
			LabelImage struct {
				LabelImageFormat struct {
					Code string
				}
				GraphicImage string
				HTMLImage    string
			}
		}
	}
}

func ShipmentConfirm(req *ShipmentConfirmRequest) (*ShipmentConfirmResponse, error) {
	buf := new(bytes.Buffer)

	if err := accessRequest(buf, req.Credentials); err != nil {
		return nil, errors.New("ups.ShipConfirm: Access request failed:\n" + err.Error())
	}

	// Load template.
	page, err := template.ParseGlob("./templates/xml/ups/*.xml")
	if err != nil {
		return nil, errors.New("ups.ShipConfirm: Unable to load template:\n" + err.Error())
	}

	// Now we can execute the template.
	err = page.ExecuteTemplate(buf, "shipmentconfirmrequest.xml", req)
	if err != nil {
		return nil, errors.New("ups.ShipConfirm: Template execution failed:\n" + err.Error())
	}

	target := "https://wwwcie.ups.com/ups.app/xml/ShipConfirm"

	rawxml, err := send(buf.Bytes(), target)
	if err != nil {
		return nil, errors.New("ups.ShipConfirm: Template execution failed:\n" + err.Error())
	}

	response := new(ShipmentConfirmResponse)

	err = xml.Unmarshal(rawxml, &response)
	if err != nil {
		return response, errors.New("ups.ShipConfirm: XML unmarshalling failed:\n" + err.Error())
	}

	return response, err
}

func NewShipmentConfirmRequest() *ShipmentConfirmRequest {
	scr := new(ShipmentConfirmRequest)
	scr.Request.RequestAction = "ShipConfirm"
	scr.Request.RequestOption = "validate"
	scr.LabelSpecification.LabelImageFormat.Code = "GIF"
	scr.LabelSpecification.LabelImageFormat.Description = "GIF"
	scr.LabelSpecification.LabelPrintMethod.Code = "GIF"
	scr.LabelSpecification.LabelPrintMethod.Description = "GIF"
	return scr
}

func ShipmentAccept(request ShipmentAcceptRequest) (*ShipmentAcceptResponse, error) {
	buf := new(bytes.Buffer)

	if err := accessRequest(buf, request.Credentials); err != nil {
		return nil, errors.New("ups.ShipmentAccept: Access request failed:\n" + err.Error())
	}

	// Load template.
	page, err := template.ParseGlob("./templates/xml/ups/*.xml")
	if err != nil {
		return nil, errors.New("ups.ShipmentAccept: Unable to load template:\n" + err.Error())
	}

	// Now we can execute the template.
	err = page.ExecuteTemplate(buf, "shipmentacceptrequest.xml", request)
	if err != nil {
		return nil, errors.New("ups.ShipmentAccept: Template execution failed:\n" + err.Error())
	}

	target := "https://wwwcie.ups.com/ups.app/xml/ShipAccept"

	rawxml, err := send(buf.Bytes(), target)
	if err != nil {
		return nil, errors.New("ups.ShipmentAccept: Template execution failed:\n" + err.Error())
	}

	response := new(ShipmentAcceptResponse)

	err = xml.Unmarshal(rawxml, &response)
	if err != nil {
		return response, errors.New("ups.ShipmentAccept: XML unmarshalling failed:\n" + err.Error())
	}

	return response, err
}

func NewShipmentAcceptRequest(scr ShipmentConfirmResponse) *ShipmentAcceptRequest {
	sar := new(ShipmentAcceptRequest)

	sar.Request.RequestAction = "ShipAccept"
	sar.Request.TransactionReference.CustomerContext = scr.Response.TransactionReference.CustomerContext
	sar.ShipmentDigest = scr.ShipmentDigest
	return sar
}

func (sar *ShipmentAcceptResponse) ExtractLabels() ([][]byte, error) {

	var images [][]byte

	for _, value := range sar.ShipmentResults.PackageResults {
		var img bytes.Buffer
		dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(value.LabelImage.GraphicImage))
		if _, err := img.ReadFrom(dec); err != nil {
			return nil, errors.New("ups.ShipmentAccept: Failed to dump GraphicImage to buffer.:\n" + err.Error())
		}
		images = append(images, img.Bytes())
	}

	return images, nil
}

//func (sar *ShipmentAcceptResponse) ExtractHTML() ([]byte, error) {
//
//	var img bytes.Buffer
//
//	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(sar.ShipmentResults.PackageResults.LabelImage.HTMLImage))
//
//	if _, err := img.ReadFrom(dec); err != nil {
//		return nil, errors.New("ups.ShipmentAccept: Failed to dump HTMLImage to buffer.:\n" + err.Error())
//	}
//
//	return img.Bytes(), nil
//}
