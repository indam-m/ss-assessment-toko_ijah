package model

// ItemValueReportItem is a struct for each item value report row
type ItemValueReportItem struct {
	SKU                  string `json:"sku"`
	Name                 string `json:"name"`
	Amount               int64  `json:"amount"`
	AveragePurchasePrice int64  `json:"average_purchase_price"`
	Total                int64  `json:"total"`
}

// SellingReportItem is a struct for each selling report row
type SellingReportItem struct {
	OrderID       string `json:"order_id"`
	Time          string `json:"time"`
	SKU           string `json:"sku"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	SellingPrice  int64  `json:"selling_price"`
	Total         int64  `json:"total"`
	PurchasePrice int64  `json:"purchase_price"`
	Profit        int64  `json:"profit"`
}

// ItemValueReport is a struct for the item value report
type ItemValueReport struct {
	PrintedDate   string                `json:"printed_date"`
	AmountOfSKU   int64                 `json:"amount_of_sku"`
	AmountOfItems int64                 `json:"amount_of_items"`
	TotalValue    int64                 `json:"total_value"`
	Rows          []ItemValueReportItem `json:"rows"`
}

// SellingReport is a struct for the selling report
type SellingReport struct {
	PrintedDate      string              `json:"printed_date"`
	Date             string              `json:"date"`
	TotalTurnover    int64               `json:"total_turnover"`
	TotalGrossProfit int64               `json:"total_gross_profit"`
	TotalSelling     int64               `json:"total_selling"`
	TotalItem        int64               `json:"total_item"`
	Rows             []SellingReportItem `json:"rows"`
}
