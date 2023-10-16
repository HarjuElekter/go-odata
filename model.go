package odata

type Model struct {
	OdataContext  string      `json:"@odata.context"`
	Value         interface{} `json:"value"`
	OdataNextLink string      `json:"@odata.nextLink,omitempty"`
}
