package models

// KeyFloatPairs struct
type KeyFloatPairs map[string]*float64

// Location struct
type Location struct {
	Country *string `json:"country"`
	City    *string `json:"city"`
}

// Address struct
type Address struct {
	Address  *string `json:"address"`
	Code     *string `json:"code"`
	City     *string `json:"city"`
	State    *string `json:"state"`
	PostCode *string `json:"postCode"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

// Contact struct
type Contact struct {
	Name           *string   `json:"name"`
	Language       *string   `json:"language"`
	NotifyMode     *string   `json:"notifyMode"`
	AttachmentType *string   `json:"attachmentType"`
	Email          *string   `json:"email"`
	Fax            *string   `json:"fax"`
	Phone          *[]string `json:"phone"`
}

// Note struct
type Note struct {
	Type      *string `json:"type"`
	Text      *string `json:"text"`
	CreatedAt *int64  `json:"createdAt"`
}

// Organisation struct
type Organisation struct {
	Name      *string   `json:"name"`
	Location  *string   `json:"location"`
	Addresses *[]string `json:"addresses"`
	Contacts  *[]string `json:"contacts"`
	EDICode   *string   `json:"ediCode"`
	OwnerCode *string   `json:"ownerCode"`
}

// DocumentLink struct
type DocumentLink struct {
	Date        *int64  `json:"date"`
	Description *string `json:"description"`
	Link        *string `json:"link"`
}
