package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemIn is used as the controller struct
type ItemIn struct{}

func getInitItemIn(r *http.Request) model.ItemIn {
	var itemIn model.ItemIn
	itemIn.ID, _ = strconv.ParseInt(r.FormValue("ID"), 10, 64)
	itemIn.SKU = r.FormValue("SKU")
	itemIn.AmountOrders, _ = strconv.ParseInt(r.FormValue("AmountOrders"), 10, 64)
	itemIn.AmountReceived, _ = strconv.ParseInt(r.FormValue("AmountReceived"), 10, 64)
	itemIn.PurchasePrice, _ = strconv.ParseInt(r.FormValue("PurchasePrice"), 10, 64)
	itemIn.ReceiptNumber = r.FormValue("ReceiptNumber")
	itemIn.Notes = r.FormValue("Notes")
	itemIn.Time = getStringDate(r.FormValue("Time"))
	return itemIn
}

func getItemInList(w http.ResponseWriter, r *http.Request) []model.ItemIn {
	rows, err := database.Query(`
		SELECT item_in.id, item_in.time, item_in.sku,
		item_amount.name, item_in.amount_orders, item_in.amount_received,
		item_in.purchase_price, item_in.receipt_number, item_in.notes,
		(item_in.purchase_price * item_in.amount_orders) AS total
		FROM item_in
		INNER JOIN item_amount ON item_in.sku=item_amount.sku
	`)
	checkInternalServerError(err, w)

	var itemIns []model.ItemIn
	var itemIn model.ItemIn

	for rows.Next() {
		err = rows.Scan(
			&itemIn.ID, &itemIn.Time,
			&itemIn.SKU, &itemIn.Name,
			&itemIn.AmountOrders, &itemIn.AmountReceived,
			&itemIn.PurchasePrice, &itemIn.ReceiptNumber,
			&itemIn.Notes, &itemIn.Total,
		)
		checkInternalServerError(err, w)
		itemIns = append(itemIns, itemIn)
	}
	return itemIns
}

// GetItemIns returns list of item_ins
func (ctrl ItemIn) GetItemIns(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	itemIns := getItemInList(w, r)

	t, err := json.Marshal(itemIns)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}

// ExportItemIns exports list of item_ins
func (ctrl ItemIn) ExportItemIns(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-in", 301)
	}
	itemIns := getItemInList(w, r)

	// creating csv file
	f, err := os.Create("catatan_barang_masuk.csv")
	checkInternalServerError(err, w)
	defer f.Close()
	csvw := csv.NewWriter(f)
	defer csvw.Flush()

	csvw.Write([]string{"ID", "Waktu", "SKU", "Nama Barang", "Jumlah Pemesanan", "Jumlah Diterima", "Harga Beli", "Total", "Nomer Kwitansi", "Catatan"})
	for _, v := range itemIns {
		err := csvw.Write([]string{
			convertToStr(v.ID),
			getDateTimeStr(v.Time),
			v.SKU,
			v.Name,
			convertToStr(v.AmountOrders),
			convertToStr(v.AmountReceived),
			convertToStr(v.PurchasePrice),
			convertToStr(v.Total),
			v.ReceiptNumber,
			v.Notes,
		})
		checkInternalServerError(err, w)
	}
	// done creating csv file

	t, err := json.Marshal(itemIns)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}

// CreateItemIn creates an item_in from request
func (ctrl ItemIn) CreateItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-in", 301)
	}
	itemIn := getInitItemIn(r)

	// Save to database
	stmt, err := database.Prepare(`
		INSERT INTO item_in(sku, time, amount_orders,
		amount_received, purchase_price, receipt_number, notes)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Fprintln(w, "Prepare query error")
		fmt.Fprintf(w, err.Error())
	}
	_, err = stmt.Exec(
		itemIn.SKU, itemIn.Time, itemIn.AmountOrders,
		itemIn.AmountReceived, itemIn.PurchasePrice,
		itemIn.ReceiptNumber, itemIn.Notes,
	)
	if err != nil {
		fmt.Fprintln(w, "Execute query error")
		fmt.Fprintf(w, err.Error())
	}
	txt, _ := json.Marshal(itemIn)
	fmt.Fprintln(w, "Creating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// UpdateItemIn updates an item_in from request using item_in ID
func (ctrl ItemIn) UpdateItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-in", 301)
	}
	itemIn := getInitItemIn(r)

	stmt, err := database.Prepare(`
		UPDATE item_in SET time=?, sku=?,
		amount_orders=?, amount_received=?,
		purchase_price=?, receipt_number=?,
		notes=?
		WHERE id=?
	`)
	checkInternalServerError(err, w)
	res, err := stmt.Exec(
		itemIn.Time, itemIn.SKU,
		itemIn.AmountOrders, itemIn.AmountReceived,
		itemIn.PurchasePrice, itemIn.ReceiptNumber,
		itemIn.Notes, itemIn.ID,
	)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	txt, _ := json.Marshal(itemIn)
	fmt.Fprintln(w, "Updating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// DeleteItemIn deletes an item_in using requested ID
func (ctrl ItemIn) DeleteItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-in", 301)
	}
	var itemID = r.FormValue("ID")
	stmt, err := database.Prepare("DELETE FROM item_in WHERE id=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	fmt.Fprintf(w, "Deleting succeeded!")
}
