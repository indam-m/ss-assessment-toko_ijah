package model

import "time"

// ItemIn Structure
type ItemIn struct {
	ID             int64     `json:"id"`
	Time           time.Time `json:"time"`
	SKU            string    `json:"sku"`
	Name           string    `jsn:"name"`
	AmountOrders   int64     `json:"amount_orders"`
	AmountReceived int64     `json:"amount_received"`
	PurchasePrice  int64     `json:"purchase_price"`
	ReceiptNumber  string    `json:"receipt_number"`
	Notes          string    `json:"notes"`
}
