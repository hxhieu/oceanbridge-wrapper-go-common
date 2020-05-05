package models

// DashboardShipment struct
type DashboardShipment struct {
	Number         *string        `json:"number"`
	Origin         *Location      `json:"origin"`
	Destination    *Location      `json:"destination"`
	DetentionSpend *KeyFloatPairs `json:"detentionSpend"`
	Prices         *KeyFloatPairs `json:"prices"`
	ETD            *int64         `json:"etd"`
	ETA            *int64         `json:"eta"`
	DeliveredDate  *int64         `json:"deliveredDate"`
}

// Dashboard struct
type Dashboard struct {
	Shipments *[]DashboardShipment
}

// Map maps the related info from the shipments
func (dest *Dashboard) Map(src *GetShipmentListResult) {
	shipments := make([]DashboardShipment, 0)
	for _, s := range *src.Shipments {
		shipment := DashboardShipment{
			Number:        s.Number,
			Prices:        s.Prices,
			ETA:           s.ETA,
			ETD:           s.ETD,
			DeliveredDate: s.DeliveredDate,
		}
		if s.Origin != nil {
			origin := (*src.Lookup)[*s.Origin].(Location)
			shipment.Origin = &origin
		}
		if s.Destination != nil {
			destination := (*src.Lookup)[*s.Destination].(Location)
			shipment.Destination = &destination
		}
		shipments = append(shipments, shipment)
	}
	dest.Shipments = &shipments
}
