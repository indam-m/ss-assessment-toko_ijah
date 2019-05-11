package model

import "time"

// ItemOut Structure
type ItemOut struct {
	ID           int64     `json:"id"`
	Time         time.Time `json:"time"`
	SKU          string    `json:"sku"`
	Name         string    `jsn:"name"`
	AmountOut    int64     `json:"amount_out"`
	SellingPrice int64     `json:"selling_price"`
	OrderID      string    `json:"order_id"`
	Notes        string    `json:"notes"`
}
