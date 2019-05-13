package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/indam-m/ss-assessment-toko_ijah/model"
)

// ItemOut is used as the controller struct
type ItemOut struct{}

func getInitItemOut(r *http.Request) model.ItemOut {
	var itemOut model.ItemOut
	itemOut.ID, _ = strconv.ParseInt(r.FormValue("ID"), 10, 64)
	itemOut.SKU = r.FormValue("SKU")
	itemOut.AmountOut, _ = strconv.ParseInt(r.FormValue("AmountOut"), 10, 64)
	itemOut.SellingPrice, _ = strconv.ParseInt(r.FormValue("SellingPrice"), 10, 64)
	itemOut.OrderID = r.FormValue("OrderID")
	itemOut.Notes = r.FormValue("Notes")
	itemOut.Time = convertDateForSQL(r.FormValue("Time"))
	return itemOut
}

func getItemOutList(w http.ResponseWriter, r *http.Request) []model.ItemOut {
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
	return itemOuts
}

// GetItemOuts returns list of item_outs
func (ctrl ItemOut) GetItemOuts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	itemOuts := getItemOutList(w, r)

	t, err := template.New("item-out.html").Funcs(getTemplateFunc()).ParseFiles("assets/item-out.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, itemOuts)
	checkInternalServerError(err, w)
}

// ExportItemOuts exports list of item_outs
func (ctrl ItemOut) ExportItemOuts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-out", 301)
	}
	itemOuts := getItemOutList(w, r)

	// creating csv file
	f, err := os.Create("catatan_barang_keluar.csv")
	checkInternalServerError(err, w)
	defer f.Close()
	csvw := csv.NewWriter(f)
	defer csvw.Flush()

	csvw.Write([]string{"ID", "Waktu", "SKU", "Nama Barang", "Jumlah Keluar", "Harga Jual", "Total", "Catatan"})
	for _, v := range itemOuts {
		err := csvw.Write([]string{
			convertToStr(v.ID),
			convertToUITime(v.Time),
			v.SKU,
			v.Name,
			convertToStr(v.AmountOut),
			convertToStr(v.SellingPrice),
			convertToStr(v.Total),
			v.Notes,
		})
		checkInternalServerError(err, w)
	}
	// done creating csv file

	fmt.Fprintln(w, exportSuccess)
}

// CreateItemOut creates an item_out from request
func (ctrl ItemOut) CreateItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-out", 301)
	}
	itemOut := getInitItemOut(r)

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
	fmt.Fprintln(w, createSuccess)
	fmt.Fprintf(w, string(txt))
}

// UpdateItemOut updates an item_out from request using item_out ID
func (ctrl ItemOut) UpdateItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-out", 301)
	}
	itemOut := getInitItemOut(r)

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
	fmt.Fprintln(w, updateSuccess)
	fmt.Fprintf(w, string(txt))
}

// DeleteItemOut deletes an item_out using requested ID
func (ctrl ItemOut) DeleteItemOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-out", 301)
	}
	var itemID = r.FormValue("ID")
	stmt, err := database.Prepare("DELETE FROM item_out WHERE id=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(itemID)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	fmt.Fprintf(w, deleteSuccess)
}

// ImportItemOuts imports item_out list from csv file
func (ctrl ItemOut) ImportItemOuts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/item-out", 301)
	}
	f, err := os.Open(r.FormValue("FileName"))
	if err != nil {
		checkInternalServerError(err, w)
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	skipHeader := false
	sqlStr := `
		INSERT INTO item_out(time, sku, amount_out,
		selling_price, order_id, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			checkInternalServerError(err, w)

		} else if skipHeader {
			re := regexp.MustCompile("ID(-\\d+)+")
			match := re.FindStringSubmatch(row[6])
			var orderID string
			if len(match) > 0 {
				orderID = match[0]
			}
			vals := []interface{}{
				convertDateForSQL(row[0]), row[1],
				convertToInt(row[3]), convertToInt(row[4]),
				orderID, row[6],
			}
			execImport(sqlStr, vals, w)
		} else {
			skipHeader = true
		}
	}
	if err == nil {
		fmt.Fprintln(w, importSuccess)
	}
}
