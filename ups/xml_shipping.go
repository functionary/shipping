package ups

import ()

type ShipmentConfirmRequest struct {
	Shipper  ShipperType
	ShipTo   ShipToType
	ShipFrom ShipFromType

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
			Code        LabelImageFormatCode
			Description string
		}
		LabelStockSize struct {
			// In inches, whole numbers only.
			Width  int
			Height int
		}
	}
	Shipment ShipmentType
	Service  struct {
		Code        ServiceCode
		Description string
	}
	Packages []PackageType `xml:"Package"`
}

type ShipmentAcceptRequest struct {
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
