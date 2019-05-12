package model

// ItemOut Structure
type ItemOut struct {
	ID           int64  `json:"id"`
	Time         string `json:"time"`
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	AmountOut    int64  `json:"amount_out"`
	SellingPrice int64  `json:"selling_price"`
	Total        int64  `json:"total"`
	OrderID      string `json:"order_id"`
	Notes        string `json:"notes"`
}
