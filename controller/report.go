package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	// "strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// Report is used as the controller struct
type Report struct{}

// GetItemValueReport returns the item value report
func (ctrl Report) GetItemValueReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := database.Query(`
		SELECT item_in.sku, item_amount.name,
		SUM (item_in.amount_received) AS amount,
		SUM (item_in.amount_received * item_in.purchase_price) AS total_price
		FROM item_in INNER JOIN item_amount ON item_in.sku=item_amount.sku
		GROUP BY item_in.sku
	`)

	var item model.ItemValueReportItem
	var report model.ItemValueReport

	for rows.Next() {
		var totalPrice int64
		err = rows.Scan(
			&item.SKU, &item.Name,
			&item.Amount, &totalPrice,
		)
		checkInternalServerError(err, w)
		item.AveragePurchasePrice = totalPrice / item.Amount
		item.Total = item.Amount * item.AveragePurchasePrice
		report.AmountOfSKU++
		report.AmountOfItems += item.Amount
		report.TotalValue += item.Total
		report.Rows = append(report.Rows, item)
	}
	report.PrintedDate = time.Now().Format("2 January 2006")
	t, err := json.Marshal(report)
	checkInternalServerError(err, w)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}

func getFilteringDate(str string, isFrom bool) string {
	theTime, err := time.Parse("2 January 2006", str)
	if err != nil {
		theTime = time.Now()
	}
	if !isFrom {
		theTime = theTime.AddDate(0, 0, 1)
	}
	return theTime.Format("2006-01-02")
}

// GetSellingReport returns the selling report
func (ctrl Report) GetSellingReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	dateFrom := r.FormValue("DateFrom")
	dateTo := r.FormValue("DateTo")

	rows, err := database.Query(`
		SELECT item_out.order_id, item_out.time, item_out.sku,
		item_amount.name, item_out.amount_out, item_out.selling_price,
		(item_out.amount_out * item_out.selling_price) AS total_price,
		(grouped_item_in.total_price / grouped_item_in.amount) AS purchase_price
		FROM item_out
		INNER JOIN item_amount ON item_out.sku=item_amount.sku
		INNER JOIN
		(SELECT item_in.sku, 
		SUM (item_in.amount_received) AS amount,
		SUM (item_in.amount_received * item_in.purchase_price) AS total_price
		FROM item_in
		GROUP BY item_in.sku) grouped_item_in
		ON item_out.sku=grouped_item_in.sku
		WHERE item_out.time >= ?
		AND item_out.time < ?
	`, getFilteringDate(dateFrom, true), getFilteringDate(dateTo, false))

	var (
		item   model.SellingReportItem
		report model.SellingReport
	)

	for rows.Next() {
		fmt.Println("aaa")
		err = rows.Scan(
			&item.OrderID, &item.Time,
			&item.SKU, &item.Name,
			&item.Amount, &item.SellingPrice,
			&item.Total, &item.PurchasePrice,
		)
		checkInternalServerError(err, w)
		item.Profit = item.Total - (item.PurchasePrice * item.Amount)

		report.TotalTurnover += item.Total
		report.TotalGrossProfit += item.Profit
		report.TotalItem += item.Amount
		if len(item.OrderID) > 0 {
			report.TotalSelling++
		}
		report.Rows = append(report.Rows, item)
	}
	report.Date = dateFrom + " - " + dateTo
	report.PrintedDate = time.Now().Format("2 January 2006")
	t, err := json.Marshal(report)
	checkInternalServerError(err, w)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}
