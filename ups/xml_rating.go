package ups

import ()

type RatingServiceSelectionResponse struct {
	Response struct {
		ResponseStatusCode        int
		ResponseStatusDescription string `xml:",omitempty"`
		TransactionReference      struct {
			CustomerContext string
			XpciVersion     string
		}
	}
	RatedShipment []struct {
		Service struct {
			Code int
		}
		RatedShipmentWarning string
		BillingWeight        struct {
			UnitOfMeasurement struct {
				Code string
			}
			Weight float64
		}
		TransportationCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}

		ServiceOptionsCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}

		GuaranteedDaysToDelivery string `xml:",omitempty"`
		ScheduledDeliveryTime    string `xml:",omitempty"`

		TotalCharges struct {
			CurrencyCode  string
			MonetaryValue float64
		}

		RatedPackage struct {
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
			Weight        float64
			BillingWeight struct {
				UnitOfMeasurement struct {
					Code string
				}
				Weight float64
			}
		}
	}
}

type RatingServiceSelectionRequest struct {
	Request struct {
		TransactionReference struct {
			CustomerContext string
			XpciVersion     string
		}
		RequestAction string
		RequestOption string
	}
	PickupType struct {
		Code        PickupTypeCode
		Description string
	}
	CustomerClassification []struct {
		Code CustomerClassificationCode
	}
	Shipment ShipmentType
}
