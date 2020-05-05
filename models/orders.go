package models

// Order struct
type Order struct {
	Shipment *Shipment `json:"shipment"`
}

// GetOrderListResult struct
type GetOrderListResult struct {
	TotalRows    uint64   `json:"totalRows"`
	ReturnedRows uint64   `json:"returnedRows"`
	Orders       *[]Order `json:"orders"`
}
