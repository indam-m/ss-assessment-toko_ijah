package model

// ItemAmount Structure
type ItemAmount struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}
