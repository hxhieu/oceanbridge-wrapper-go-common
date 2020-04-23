package models

import (
	"time"
)

// Shipment struct
type Shipment struct {
	Number           string    `json:"id" example:"NMTOKEN"`
	HouseBill        string    `json:"houseBill" example:"NMTOKEN"`
	Shipper          string    `json:"shipper"`
	Consignee        string    `json:"consignee"`
	ClientReference  string    `json:"clientRef" example:"normalizedString"`
	GoodsDescription string    `json:"goodsDesc" example:"normalizedString"`
	ServiceLevel     string    `json:"serviceLevel" example:"NMTOKEN"`
	Origin           string    `json:"origin"`
	Destination      string    `json:"destination"`
	ETD              time.Time `json:"etd"`
	ETA              time.Time `json:"eta"`
	DeliveredDate    time.Time `json:"deliveredDate"`
	Size             float32   `json:"size"`
	Weight           float32   `json:"weight"`
	Quantity         int32     `json:"quantity"`
}

// Order struct
type Order struct {
	Shipment      Shipment `json:"shipment"`
	DocumentLinks []string `json:"documentLinks"`
}

// GetShipmentListResult struct
type GetShipmentListResult struct {
	TotalRows    uint32     `json:"totalRows"`
	ReturnedRows uint32     `json:"returnedRows"`
	Shipments    []Shipment `json:"shipments"`
}

// GetOrderListResult struct
type GetOrderListResult struct {
	TotalRows    uint32  `json:"totalRows"`
	ReturnedRows uint32  `json:"returnedRows"`
	Orders       []Order `json:"orders"`
}
