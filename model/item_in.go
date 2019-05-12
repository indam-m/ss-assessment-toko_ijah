package model

// ItemIn Structure
type ItemIn struct {
	ID             int64  `json:"id"`
	Time           string `json:"time"`
	SKU            string `json:"sku"`
	Name           string `json:"name"`
	AmountOrders   int64  `json:"amount_orders"`
	AmountReceived int64  `json:"amount_received"`
	PurchasePrice  int64  `json:"purchase_price"`
	Total          int64  `json:"total"`
	ReceiptNumber  string `json:"receipt_number"`
	Notes          string `json:"notes"`
}
