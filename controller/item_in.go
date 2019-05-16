package controller

import (
	"encoding/csv"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemIn is used as the controller struct
type ItemIn struct{}

var itemInHome = "/item-in"

func getInitItemIn(r *http.Request) model.ItemIn {
	var itemIn model.ItemIn
	itemIn.ID, _ = strconv.ParseInt(r.FormValue("ID"), 10, 64)
	itemIn.SKU = r.FormValue("SKU")
	itemIn.AmountOrders, _ = strconv.ParseInt(r.FormValue("AmountOrders"), 10, 64)
	itemIn.AmountReceived, _ = strconv.ParseInt(r.FormValue("AmountReceived"), 10, 64)
	itemIn.PurchasePrice, _ = strconv.ParseInt(r.FormValue("PurchasePrice"), 10, 64)
	itemIn.ReceiptNumber = r.FormValue("ReceiptNumber")
	itemIn.Notes = r.FormValue("Notes")
	itemIn.Time = convertDateForSQL(r.FormValue("Time"))
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
	alertFromCookie(w, r)

	itemIns := getItemInList(w, r)

	t, err := template.New("item-in.html").Funcs(getTemplateFunc()).ParseFiles("assets/item-in.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, itemIns)
	checkInternalServerError(err, w)
}

// ExportItemIns exports list of item_ins
func (ctrl ItemIn) ExportItemIns(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemInHome, "")
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
			convertToUITime(v.Time),
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
	redirectWithAlert(w, r, itemInHome, exportSuccess)
}

// CreateItemIn creates an item_in from request
func (ctrl ItemIn) CreateItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemInHome, "")
	}
	itemIn := getInitItemIn(r)

	// Save to database
	stmt, err := database.Prepare(`
		INSERT INTO item_in(sku, time, amount_orders,
		amount_received, purchase_price, receipt_number, notes)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		redirectWithAlert(w, r, itemAmountHome, "Prepare query error: "+err.Error())
	}
	_, err = stmt.Exec(
		itemIn.SKU, itemIn.Time, itemIn.AmountOrders,
		itemIn.AmountReceived, itemIn.PurchasePrice,
		itemIn.ReceiptNumber, itemIn.Notes,
	)
	if err != nil {
		redirectWithAlert(w, r, itemAmountHome, "Execute query error: "+err.Error())
	}
	redirectWithAlert(w, r, itemInHome, createSuccess)
}

// UpdateItemIn updates an item_in from request using item_in ID
func (ctrl ItemIn) UpdateItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemInHome, "")
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
	redirectWithAlert(w, r, itemInHome, updateSuccess)
}

// DeleteItemIn deletes an item_in using requested ID
func (ctrl ItemIn) DeleteItemIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemInHome, "")
	}
	var itemID = r.FormValue("ID")
	stmt, err := database.Prepare("DELETE FROM item_in WHERE id=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	redirectWithAlert(w, r, itemInHome, deleteSuccess)
}

// ImportItemIns imports item_in list from csv file
func (ctrl ItemIn) ImportItemIns(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemInHome, "")
	}
	f, err := os.Open(r.FormValue("FileName"))
	if err != nil {
		checkInternalServerError(err, w)
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	skipHeader := false
	sqlStr := `INSERT INTO item_in(time, sku, amount_orders,
		amount_received, purchase_price, receipt_number, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			checkInternalServerError(err, w)

		} else if skipHeader {
			vals := []interface{}{
				convertDateForSQL(row[0]), row[1],
				convertToInt(row[3]), convertToInt(row[4]),
				convertToInt(row[5]), row[7], row[8],
			}
			execImport(sqlStr, vals, w)
		} else {
			skipHeader = true
		}
	}
	if err == nil {
		redirectWithAlert(w, r, itemInHome, importSuccess)
	}
}
