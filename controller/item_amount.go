package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemAmount is used as the controller struct
type ItemAmount struct{}

const itemAmountHome = "/item-amount"

func getInitItemAmount(r *http.Request) model.ItemAmount {
	var itemAmount model.ItemAmount
	itemAmount.SKU = r.FormValue("SKU")
	itemAmount.Name = r.FormValue("Name")
	itemAmount.Quantity, _ = strconv.ParseInt(r.FormValue("Quantity"), 10, 64)
	return itemAmount
}

func getItemAmountList(w http.ResponseWriter, r *http.Request) []model.ItemAmount {
	rows, err := database.Query("SELECT * FROM item_amount")
	checkInternalServerError(err, w)

	var itemAmounts []model.ItemAmount
	var itemAmount model.ItemAmount
	for rows.Next() {
		err = rows.Scan(&itemAmount.SKU, &itemAmount.Name, &itemAmount.Quantity)
		checkInternalServerError(err, w)
		itemAmounts = append(itemAmounts, itemAmount)
	}
	return itemAmounts
}

// GetItemAmounts returns list of item_amounts
func (ctrl ItemAmount) GetItemAmounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	alertFromCookie(w, r)

	itemAmounts := getItemAmountList(w, r)
	t, err := template.New("item-amount.html").ParseFiles("assets/item-amount.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, itemAmounts)
	checkInternalServerError(err, w)
}

// ExportItemAmounts exports list of item_amounts
func (ctrl ItemAmount) ExportItemAmounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemAmountHome, "")
	}
	itemAmounts := getItemAmountList(w, r)

	// creating csv file
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=catatan_jumlah_barang.csv")
	csvw := csv.NewWriter(w)
	defer csvw.Flush()

	csvw.Write([]string{"SKU", "Nama Item", "Jumlah Sekarang"})
	for _, v := range itemAmounts {
		err := csvw.Write([]string{
			v.SKU,
			v.Name,
			convertToStr(v.Quantity),
		})
		checkInternalServerError(err, w)
	}
	// done creating csv file
}

// GetItemAmount returns an item_amounts based on SKU
func (ctrl ItemAmount) GetItemAmount(w http.ResponseWriter, r *http.Request, itemSKU string) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := database.Query("SELECT * FROM item_amount WHERE sku='" + itemSKU + "'")
	checkInternalServerError(err, w)

	var itemAmount model.ItemAmount
	for rows.Next() {
		err = rows.Scan(&itemAmount.SKU, &itemAmount.Name, &itemAmount.Quantity)
		checkInternalServerError(err, w)
	}
	t, err := json.Marshal(itemAmount)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}

// CreateItemAmount creates an item_amount from request
func (ctrl ItemAmount) CreateItemAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemAmountHome, "")
	}
	itemAmount := getInitItemAmount(r)

	// Save to database
	stmt, err := database.Prepare(`
		INSERT INTO item_amount(sku, name, quantity)
		VALUES(?, ?, ?)
	`)
	if err != nil {
		redirectWithAlert(w, r, itemAmountHome, "Prepare query error: "+err.Error())
		return
	}
	_, err = stmt.Exec(itemAmount.SKU, itemAmount.Name, itemAmount.Quantity)
	if err != nil {
		redirectWithAlert(w, r, itemAmountHome, "Execute query error: "+err.Error())
		return
	}
	redirectWithAlert(w, r, itemAmountHome, createSuccess)
}

// UpdateItemAmount updates an item_amount from request
func (ctrl ItemAmount) UpdateItemAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemAmountHome, "")
	}
	itemAmount := getInitItemAmount(r)

	stmt, err := database.Prepare(`
		UPDATE item_amount SET name=?, quantity=?
		WHERE sku=?
	`)
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemAmount.Name, itemAmount.Quantity, itemAmount.SKU)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	redirectWithAlert(w, r, itemAmountHome, updateSuccess)
}

// DeleteItemAmount deletes an item_amount using requested SKU
func (ctrl ItemAmount) DeleteItemAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemAmountHome, "")
	}
	var itemSKU = r.FormValue("SKU")
	stmt, err := database.Prepare("DELETE FROM item_amount WHERE sku=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemSKU)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	redirectWithAlert(w, r, itemAmountHome, deleteSuccess)
}

// ImportItemAmounts imports item_amount list from csv file
func (ctrl ItemAmount) ImportItemAmounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectWithAlert(w, r, itemAmountHome, "")
	}
	f, err := os.Open(r.FormValue("FileName"))
	if err != nil {
		checkInternalServerError(err, w)
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	skipHeader := false
	sqlStr := "INSERT OR REPLACE INTO item_amount(sku, name, quantity) VALUES(?, ?, ?)"
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
				row[0], row[1], convertToInt(row[2]),
			}
			err = execImport(sqlStr, vals, w)
		} else {
			skipHeader = true
		}
	}
	if err == nil {
		redirectWithAlert(w, r, itemAmountHome, importSuccess)
	}
}
