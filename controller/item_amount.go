package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemAmount is used as the controller struct
type ItemAmount struct{}

// GetItemAmounts returns list of item_amounts
func (ctrl ItemAmount) GetItemAmounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := database.Query("SELECT * FROM item_amount")
	checkInternalServerError(err, w)

	var itemAmounts []model.ItemAmount
	var itemAmount model.ItemAmount
	for rows.Next() {
		err = rows.Scan(&itemAmount.SKU, &itemAmount.Name, &itemAmount.Quantity)
		checkInternalServerError(err, w)
		itemAmounts = append(itemAmounts, itemAmount)
	}
	t, err := json.Marshal(itemAmounts)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
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
		http.Redirect(w, r, "/", 301)
	}
	var itemAmount model.ItemAmount
	itemAmount.SKU = r.FormValue("SKU")
	itemAmount.Name = r.FormValue("Name")
	itemAmount.Quantity, _ = strconv.ParseInt(r.FormValue("Quantity"), 10, 64)

	// Save to database
	stmt, err := database.Prepare(`
		INSERT INTO item_amount(sku, name, quantity)
		VALUES(?, ?, ?)
	`)
	if err != nil {
		fmt.Fprintln(w, "Prepare query error")
		fmt.Fprintf(w, err.Error())
	}
	_, err = stmt.Exec(itemAmount.SKU, itemAmount.Name, itemAmount.Quantity)
	if err != nil {
		fmt.Fprintln(w, "Execute query error")
		fmt.Fprintf(w, err.Error())
	}
	txt, _ := json.Marshal(itemAmount)
	fmt.Fprintln(w, "Creating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// UpdateItemAmount updates an item_amount from request
func (ctrl ItemAmount) UpdateItemAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var itemAmount model.ItemAmount
	itemAmount.SKU = r.FormValue("SKU")
	itemAmount.Name = r.FormValue("Name")
	itemAmount.Quantity, _ = strconv.ParseInt(r.FormValue("Quantity"), 10, 64)

	stmt, err := database.Prepare(`
		UPDATE item_amount SET name=?, quantity=?
		WHERE sku=?
	`)
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemAmount.Name, itemAmount.Quantity, itemAmount.SKU)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	txt, _ := json.Marshal(itemAmount)
	fmt.Fprintln(w, "Updating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// DeleteItemAmount deletes an item_amount using requested SKU
func (ctrl ItemAmount) DeleteItemAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var itemSKU = r.FormValue("SKU")
	stmt, err := database.Prepare("DELETE FROM item_amount WHERE sku=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemSKU)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	fmt.Fprintf(w, "Deleting succeeded!")
}