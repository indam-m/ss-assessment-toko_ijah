package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemOut is used as the controller struct
type ItemOut struct{}

func getItemOut(r *http.Request) model.ItemOut {
	var itemOut model.ItemOut
	itemOut.ID, _ = strconv.ParseInt(r.FormValue("ID"), 10, 64)
	itemOut.SKU = r.FormValue("SKU")
	itemOut.AmountOut, _ = strconv.ParseInt(r.FormValue("AmountOut"), 10, 64)
	itemOut.SellingPrice, _ = strconv.ParseInt(r.FormValue("SellingPrice"), 10, 64)
	itemOut.OrderID = r.FormValue("OrderID")
	itemOut.Notes = r.FormValue("Notes")
	itemOut.Time = getStringDate(r.FormValue("Time"))
	return itemOut
}

// GetItemOuts returns list of item_outs
func (ctrl ItemOut) GetItemOuts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := database.Query(`
		SELECT item_out.id, item_out.time, item_out.sku,
		item_amount.name, item_out.amount_out, item_out.selling_price,
		item_out.order_id, item_out.notes
		FROM item_out
		INNER JOIN item_amount ON item_out.sku=item_amount.sku
	`)
	checkInternalServerError(err, w)

	var itemOuts []model.ItemOut
	var itemOut model.ItemOut

	for rows.Next() {
		err = rows.Scan(
			&itemOut.ID, &itemOut.Time,
			&itemOut.SKU, &itemOut.Name,
			&itemOut.AmountOut, &itemOut.SellingPrice,
			&itemOut.OrderID, &itemOut.Notes,
		)
		checkInternalServerError(err, w)
		itemOut.Total = itemOut.AmountOut * itemOut.SellingPrice
		itemOuts = append(itemOuts, itemOut)
	}
	t, err := json.Marshal(itemOuts)
	checkInternalServerError(err, w)
	fmt.Fprintf(w, string(t))
}

// CreateItemOut creates an item_out from request
func (ctrl ItemOut) CreateItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	itemOut := getItemOut(r)

	// Save to database
	stmt, err := database.Prepare(`
		INSERT INTO item_out(sku, time, amount_out,
		selling_price, order_id, notes)
		VALUES(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Fprintln(w, "Prepare query error")
		fmt.Fprintf(w, err.Error())
	}
	_, err = stmt.Exec(
		itemOut.SKU, itemOut.Time,
		itemOut.AmountOut, itemOut.SellingPrice,
		itemOut.OrderID, itemOut.Notes,
	)
	if err != nil {
		fmt.Fprintln(w, "Execute query error")
		fmt.Fprintf(w, err.Error())
	}
	txt, _ := json.Marshal(itemOut)
	fmt.Fprintln(w, "Creating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// UpdateItemOut updates an item_out from request using item_out ID
func (ctrl ItemOut) UpdateItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	itemOut := getItemOut(r)

	stmt, err := database.Prepare(`
		UPDATE item_out SET time=?, sku=?,
		amount_out=?, selling_price=?,
		order_id=?, notes=?
		WHERE id=?
	`)
	checkInternalServerError(err, w)
	res, err := stmt.Exec(
		itemOut.Time, itemOut.SKU,
		itemOut.AmountOut, itemOut.SellingPrice,
		itemOut.OrderID, itemOut.Notes,
		itemOut.ID,
	)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	txt, _ := json.Marshal(itemOut)
	fmt.Fprintln(w, "Updating succeeded!")
	fmt.Fprintf(w, string(txt))
}

// DeleteItemOut deletes an item_out using requested ID
func (ctrl ItemOut) DeleteItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var itemID = r.FormValue("ID")
	stmt, err := database.Prepare("DELETE FROM item_out WHERE id=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	fmt.Fprintf(w, "Deleting succeeded!")
}
