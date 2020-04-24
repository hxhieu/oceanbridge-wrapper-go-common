package shipmentservices

// OceanBridgeAuthHeader to add auth header to the envelop body
type OceanBridgeAuthHeader struct {
	Header    *WebTrackerSOAPHeader `xml:"WebTrackerSOAPHeader,omitempty" json:"WebTrackerSOAPHeader,omitempty" yaml:"WebTrackerSOAPHeader,omitempty"`
	Namespace string                `xml:"xmlns,attr,omitempty"`
}
