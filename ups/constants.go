package ups

import ()

/*

Package Type

If no container is specified, RAVE assumes UPS Package, i.e., type 02.
	01 = UPS Letter,
	02 = Customer Supplied Package,
	03 = Tube,
	04 = PAK,
	21 = UPS Express Box,
	24 = UPS 25KG Box,
	25 = UPS 10KG Box,
	30 = Pallet,
	2a = Small Express Box,
	2b = Medium Express Box,
	2c = Large Express Box,
	56 = Flats,
	57 = Parcels,
	58 = BPM,
	59 = First Class,
	60 = Priority,
	61 = Machinables,
	62 = Irregulars,
	63 = Parcel Post,
	64 = BPM Parcel,
	65 = Media Mail,
	66 = BMP Flat,
	67 = Standard Flat

Package type 24, or 25 is only allowed for shipment without return service
Packaging type must be valid for all the following:

	ShipTo country,
	ShipFrom country,
	A shipment going from ShipTo country to ShipFrom country,
	All accessorial at both the shipment and package level, and
	The shipment service type.

UPS will not accept raw wood pallets and please refer the UPS packaging
guidelines for pallets on UPS.com.

*/
type PackagingTypeCode string

const (
	PackagingTypeUnknown          PackagingTypeCode = "00"
	PackagingTypeUPSLetter        PackagingTypeCode = "01"
	PackagingTypePackage          PackagingTypeCode = "02"
	PackagingTypeTube             PackagingTypeCode = "03"
	PackagingTypePak              PackagingTypeCode = "04"
	PackagingTypeExpressBox       PackagingTypeCode = "21"
	PackagingType25KgBox          PackagingTypeCode = "24"
	PackagingType10KgBox          PackagingTypeCode = "25"
	PackagingTypePallet           PackagingTypeCode = "30"
	PackagingTypeSmallExpressBox  PackagingTypeCode = "2a"
	PackagingTypeMediumExpressBox PackagingTypeCode = "2b"
	PackagingTypeLargeExpressBox  PackagingTypeCode = "2c"
)

// This data is optional, but enhances user-friendliness.
var packageTypeNames map[PackagingTypeCode]string = map[PackagingTypeCode]string{
	PackagingTypeUnknown:          "Unknown Page",
	PackagingTypePackage:          "Package",
	PackagingTypeUPSLetter:        "UPS Letter",
	PackagingTypeTube:             "Tube",
	PackagingTypePak:              "Pak",
	PackagingTypeExpressBox:       "Express Box",
	PackagingType25KgBox:          "25KG Box",
	PackagingType10KgBox:          "10KG Box",
	PackagingTypePallet:           "Pallet",
	PackagingTypeSmallExpressBox:  "Small Express Box",
	PackagingTypeMediumExpressBox: "Medium Express Box",
	PackagingTypeLargeExpressBox:  "Large Express Box"}

/*
PICKUP TYPE:
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
type PickupTypeCode string

const (
	PickupTypeDaily            PickupTypeCode = "01"
	PickupTypeCustomerCounter  PickupTypeCode = "03"
	PickupTypeOneTime          PickupTypeCode = "06"
	PickupTypeOnCallAir        PickupTypeCode = "07"
	PickupTypeLetterCenter     PickupTypeCode = "19"
	PickupTypeAirServiceCenter PickupTypeCode = "20"
)

// This data is optional, but enhances user-friendliness.
var pickupTypeNames map[PickupTypeCode]string = map[PickupTypeCode]string{
	PickupTypeDaily:            "Daily Pickup",
	PickupTypeCustomerCounter:  "Customer Counter",
	PickupTypeOneTime:          "One Time Pickup",
	PickupTypeOnCallAir:        "On Call Air",
	PickupTypeLetterCenter:     "Letter Center",
	PickupTypeAirServiceCenter: "Air Service Center",
}

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

type ServiceCode string

const (
	// Domestic Services
	ServiceUSNextDayAirAM    ServiceCode = "14"
	ServiceUSNextDayAir      ServiceCode = "01"
	ServiceUSNextDayAirSaver ServiceCode = "13"
	ServiceUS2ndDayAirAM     ServiceCode = "59"
	ServiceUS2ndDayAir       ServiceCode = "02"
	ServiceUS3DaySelect      ServiceCode = "12"
	ServiceUSGround          ServiceCode = "03"

	// Internation Services
	ServiceIntlStandard         ServiceCode = "11"
	ServiceWorldwideExpress     ServiceCode = "07"
	ServiceWorldwideExpressPlus ServiceCode = "54"
	ServiceWorldwideExpedited   ServiceCode = "08"
	ServiceIntlSaver            ServiceCode = "65"
)

// Can this be made to 'const'?
var serviceNames map[ServiceCode]string = map[ServiceCode]string{
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

type LabelImageFormatCode string

type UnitOfMeasurementCode string

type CustomerClassificationCode string
