package models

import (
	"fmt"
	"time"

	"github.com/hxhieu/oceanbridge-wrapper-go-common/shipmentservices"
	"github.com/mitchellh/hashstructure"
)

// LookupDataPrefix map
type LookupDataPrefix struct {
	Organisation string
	Location     string
	Address      string
	Contact      string
	Documents    string
	Notes        string
	Packing      string
}

// LookDataPrefix defines the prefixes for the lookup keys
var LookDataPrefix = LookupDataPrefix{
	Organisation: "ORG",
	Location:     "LOC",
	Address:      "ADDR",
	Contact:      "CONTACT",
	Documents:    "DOC",
	Notes:        "NOTE",
	Packing:      "PKG",
}

// LookupData defines a normalize set of reference data,
// that each record can be belonged to one or many shipments
type LookupData map[string]interface{}

// Packing struct
type Packing struct {
	Type            *string        `json:"type"`
	LinePrice       *KeyFloatPairs `json:"linePrice"`
	Weight          *KeyFloatPairs `json:"weight"`
	Volume          *KeyFloatPairs `json:"volume"`
	Description     *string        `json:"description"`
	ContainerNumber *string        `json:"containerNo"`
}

// Shipment struct
type Shipment struct {
	Number           *string        `json:"number"`
	HouseBill        *string        `json:"houseBill"`
	Shipper          *string        `json:"shipper"`
	Consignee        *string        `json:"consignee"`
	GoodsDescription *string        `json:"goodsDesc"`
	ServiceLevel     *string        `json:"serviceLevel"`
	Origin           *string        `json:"origin"`
	Destination      *string        `json:"destination"`
	ETD              *int64         `json:"etd"`
	ETA              *int64         `json:"eta"`
	DeliveredDate    *int64         `json:"deliveredDate"`
	Packings         *[]string      `json:"packings"`
	Notes            *[]string      `json:"notes"`
	Documents        *[]string      `json:"documents"`
	Prices           *KeyFloatPairs `json:"prices"`
}

// GetShipmentListResult struct
type GetShipmentListResult struct {
	TotalRows    uint64      `json:"totalRows"`
	ReturnedRows uint64      `json:"returnedRows"`
	Shipments    *[]Shipment `json:"shipments"`
	Lookup       *LookupData `json:"lookup"`
}

func dateToEpoch(d *shipmentservices.DateTime, tz string) int64 {
	// From SOAP -> 2018-07-01T00:00:00
	var unix int64 = 0
	if d != nil {
		// TODO: Parse as server time
		layout := time.RFC3339
		if date, err := time.Parse(layout, fmt.Sprintf("%v%v", *d, tz)); err == nil {
			unix = date.Unix()
		}
	}
	return unix
}

func (lookup *LookupData) appendLookup(value interface{}, prefix string) *string {
	// Ignore all
	if lookup == nil {
		return nil
	}
	if hash, err := hashstructure.Hash(value, nil); err == nil {
		hashStr := fmt.Sprintf("%v_%v", prefix, hash)
		(*lookup)[hashStr] = value
		return &hashStr
	}
	return nil
}

// Map from SOAP OrgAddress to wrapper Address
func (dest *Address) Map(src *shipmentservices.OrgAddress) {
	address := ""
	if src.AddressLine1 != nil {
		address = fmt.Sprintf("%v%v", address, *src.AddressLine1)
	}
	if src.AddressLine2 != nil {
		address = fmt.Sprintf("%v, %v", address, *src.AddressLine2)
	}
	dest.Address = &address

	dest.City = src.CityOrSuburb
	dest.State = src.StateOrProvince
	dest.Code = src.AddressCode
	dest.Email = src.Email
	dest.PostCode = src.PostCode

	if src.TelephoneNumbers != nil {
		phone := ""
		for _, p := range src.TelephoneNumbers.TelephoneNumber {
			if p.Content != nil {
				phone = fmt.Sprintf("%v %v", phone, *p.Content)
			}
		}
		if len(phone) > 0 {
			dest.Phone = &phone
		}
	}
}

// Map from SOAP UNLOCO to wrapper Location
func (dest *Location) Map(src *shipmentservices.UNLOCO) {
	if len(src.City) > 0 {
		dest.City = &src.City
	}
	if len(src.Country) > 0 {
		dest.Country = &src.Country
	}
}

// Map from SOAP OrgContact to wrapper Contact
func (dest *Contact) Map(src *shipmentservices.OrgContact, lookup *LookupData) {
	dest.Name = src.Name
	dest.AttachmentType = src.AttachmentType
	dest.Language = src.Language
	dest.NotifyMode = src.NotifyMode
	dest.Fax = src.Fax
	dest.Email = src.EmailAddress
	phones := make([]string, 0)
	if src.Phone != nil {
		phones = append(phones, *src.Phone)
	}
	if src.Mobile != nil {
		phones = append(phones, *src.Mobile)
	}
	if src.HomePhone != nil {
		phones = append(phones, *src.HomePhone)
	}
	if src.OtherPhone != nil {
		phones = append(phones, *src.OtherPhone)
	}
	// Make a convention, phone ext comes last in the slice
	if src.PhoneExtension != nil {
		phones = append(phones, *src.PhoneExtension)
	}
	if len(phones) > 0 {
		dest.Phone = &phones
	}
}

// Map from SOAP Organisation to wrapper Organisation
func (dest *Organisation) Map(src *shipmentservices.Organisation, lookup *LookupData) {
	if src.OrganisationDetails != nil {
		details := *src.OrganisationDetails
		dest.Name = details.Name
		// Addresses
		if details.Addresses != nil && details.Addresses.Address != nil {
			addresses := make([]string, 0)
			for _, addr := range details.Addresses.Address {
				address := Address{}
				address.Map(addr)
				if hash := lookup.appendLookup(address, LookDataPrefix.Address); hash != nil {
					addresses = append(addresses, *hash)
				}
			}
			dest.Addresses = &addresses
		}
		// Location
		if details.Location != nil {
			location := Location{}
			location.Map(details.Location)
			dest.Location = lookup.appendLookup(location, LookDataPrefix.Location)
		}
		// Contacts
		if details.Contacts != nil && details.Contacts.Contact != nil {
			contacts := make([]string, 0)
			for _, c := range details.Contacts.Contact {
				contact := Contact{}
				contact.Map(c, lookup)
				if hash := lookup.appendLookup(contact, LookDataPrefix.Contact); hash != nil {
					contacts = append(contacts, *hash)
				}
			}
			dest.Contacts = &contacts
		}
	}
	if len(src.EDICode) > 0 {
		dest.EDICode = &src.EDICode
	}
	if len(src.OwnerCode) > 0 {
		dest.OwnerCode = &src.OwnerCode
	}
}

// Map maps a SOAP NotesNote to the wrapper Note
func (dest *Note) Map(src *shipmentservices.NotesNote) {
	typ := fmt.Sprintf("%v", src.NoteType)
	dest.Type = &typ
	dest.Text = src.NoteData
	if unix := dateToEpoch(src.NoteCreatedDateTime, ""); unix > 0 {
		dest.CreatedAt = &unix
	}
}

// Map maps a SOAP DocumentLink to the wrapper DocumentLink
func (dest *DocumentLink) Map(src *shipmentservices.DocumentLink) {
	dest.Description = src.Description
	dest.Link = src.Link
	if unix := dateToEpoch(src.Date, ""); unix > 0 {
		dest.Date = &unix
	}
}

// MapDimension maps a SOAP DimensionValue to the wrapper KeyFloatPairs
func (dest *KeyFloatPairs) MapDimension(src shipmentservices.DimensionValue) {
	key := src.DimensionType
	var value float64 = 0
	if src.Content != nil {
		value = *src.Content
	}
	(*dest)[key] = &value
}

// MapFinancial maps a SOAP FinancialValue to the wrapper KeyFloatPairs
func (dest *KeyFloatPairs) MapFinancial(src shipmentservices.FinancialValue) {
	key := src.CurrencyCode
	var value float64 = 0
	if src.Content != nil {
		value = *src.Content
	}
	(*dest)[key] = &value
}

// Map maps a SOAP WebPacking to the wrapper Packing
func (dest *Packing) Map(src *shipmentservices.WebPacking, prices *KeyFloatPairs) {
	dest.Description = src.Description
	dest.Type = src.PackType
	dest.ContainerNumber = src.ContainerNumber
	// Volume
	if src.Volume != nil {
		dest.Volume = &KeyFloatPairs{}
		dest.Volume.MapDimension(*src.Volume)
	}
	// Weight
	if src.Weight != nil {
		dest.Weight = &KeyFloatPairs{}
		dest.Weight.MapDimension(*src.Weight)
	}
	// Prices
	if src.LinePrice != nil {
		dest.LinePrice = &KeyFloatPairs{}
		dest.LinePrice.MapFinancial(*src.LinePrice)
		// Calc total value
		for k, v := range *dest.LinePrice {
			var newTotal float64 = 0
			if v != nil {
				newTotal += *v
			}
			if total, ok := (*prices)[k]; ok && total != nil {
				newTotal += *total
			}
			// Update total value
			(*prices)[k] = &newTotal
		}
	}
}

// Map maps a SOAP WebShipment to the wrapper Shipment
func (dest *Shipment) Map(src *shipmentservices.WebShipment, lookup *LookupData) {
	dest.Number = src.Number
	dest.HouseBill = src.HouseBill

	if src.Shipper != nil {
		shipper := Organisation{}
		shipper.Map(src.Shipper, lookup)
		if hash := lookup.appendLookup(shipper, LookDataPrefix.Organisation); hash != nil {
			dest.Shipper = hash
		}
	}
	if src.Consignee != nil {
		consignee := Organisation{}
		consignee.Map(src.Consignee, lookup)
		if hash := lookup.appendLookup(consignee, LookDataPrefix.Organisation); hash != nil {
			dest.Consignee = hash
		}
	}
	dest.GoodsDescription = src.GoodsDescription
	dest.ServiceLevel = src.ServiceLevel

	if src.Origin != nil {
		origin := Location{}
		origin.Map(src.Origin)
		if hash := lookup.appendLookup(origin, LookDataPrefix.Location); hash != nil {
			dest.Origin = hash
		}
	}

	if src.Destination != nil {
		destination := Location{}
		destination.Map(src.Destination)
		if hash := lookup.appendLookup(destination, LookDataPrefix.Location); hash != nil {
			dest.Destination = hash
		}
	}

	if unix := dateToEpoch(src.ETA, "Z"); unix > 0 {
		dest.ETA = &unix
	}

	if unix := dateToEpoch(src.ETD, "Z"); unix > 0 {
		dest.ETD = &unix
	}

	if unix := dateToEpoch(src.DeliveredDate, "Z"); unix > 0 {
		dest.DeliveredDate = &unix
	}
	// Notes
	if src.Notes != nil && src.Notes.Note != nil {
		notes := make([]string, 0)
		for _, n := range src.Notes.Note {
			note := Note{}
			note.Map(n)
			if hash := lookup.appendLookup(note, LookDataPrefix.Notes); hash != nil {
				notes = append(notes, *hash)
			}
		}
		dest.Notes = &notes
	}
	// Documents
	if src.DocumentLinks != nil && src.DocumentLinks.DocumentLink != nil {
		documents := make([]string, 0)
		for _, doc := range src.DocumentLinks.DocumentLink {
			document := DocumentLink{}
			document.Map(doc)
			if hash := lookup.appendLookup(document, LookDataPrefix.Documents); hash != nil {
				documents = append(documents, *hash)
			}
		}
		dest.Documents = &documents
	}
	// Packings
	if src.Packings != nil && src.Packings.Packing != nil {
		dest.Prices = &KeyFloatPairs{}
		packings := make([]string, 0)
		for _, pkg := range src.Packings.Packing {
			packing := Packing{}
			// Map also calc the total price
			packing.Map(pkg, dest.Prices)
			if hash := lookup.appendLookup(packing, LookDataPrefix.Packing); hash != nil {
				packings = append(packings, *hash)
			}
		}
		dest.Packings = &packings
	}
}

// Map maps the SOAP WebShipments to the wrapper GetShipmentListResult
func (dest *GetShipmentListResult) Map(src *shipmentservices.WebShipments, noLookup bool) {
	dest.TotalRows = src.TotalRows
	dest.ReturnedRows = src.ReturnedRows
	if !noLookup {
		dest.Lookup = &LookupData{}
	}
	// Map shipments
	shipments := make([]Shipment, 0)
	for _, src := range src.WebShipment {
		shipment := Shipment{}
		shipment.Map(src, dest.Lookup)
		shipments = append(shipments, shipment)
	}
	dest.Shipments = &shipments
}
