package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/indam-m/ss-assessment-toko_ijah/controller"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("./assets/src/dist/")))

	return r
}

func getSKU(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["sku"]
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9876"
	}
	return ":" + port
}

func main() {
	controller.Open()
	// Declare a new router
	r := mux.NewRouter()

	itemAmountCtrl := &controller.ItemAmount{}
	itemInCtrl := &controller.ItemIn{}
	itemOutCtrl := &controller.ItemOut{}
	reportCtrl := &controller.Report{}

	// declaring routers
	r.HandleFunc("/", controller.GetHome).Methods("GET")
	// item amount
	r.HandleFunc("/item-amount", itemAmountCtrl.GetItemAmounts).Methods("GET")
	r.HandleFunc("/item-amount/{sku}", func(w http.ResponseWriter, r *http.Request) {
		itemAmountCtrl.GetItemAmount(w, r, getSKU(r))
	}).Methods("GET")
	r.HandleFunc("/item-amount/create", itemAmountCtrl.CreateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/update", itemAmountCtrl.UpdateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/delete", itemAmountCtrl.DeleteItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/export", itemAmountCtrl.ExportItemAmounts).Methods("POST")
	r.HandleFunc("/item-amount/import", itemAmountCtrl.ImportItemAmounts).Methods("POST")
	// item in
	r.HandleFunc("/item-in", itemInCtrl.GetItemIns).Methods("GET")
	r.HandleFunc("/item-in/create", itemInCtrl.CreateItemIn).Methods("POST")
	r.HandleFunc("/item-in/update", itemInCtrl.UpdateItemIn).Methods("POST")
	r.HandleFunc("/item-in/delete", itemInCtrl.DeleteItemIn).Methods("POST")
	r.HandleFunc("/item-in/export", itemInCtrl.ExportItemIns).Methods("POST")
	r.HandleFunc("/item-in/import", itemInCtrl.ImportItemIns).Methods("POST")
	// item out
	r.HandleFunc("/item-out", itemOutCtrl.GetItemOuts).Methods("GET")
	r.HandleFunc("/item-out/create", itemOutCtrl.CreateItemOut).Methods("POST")
	r.HandleFunc("/item-out/update", itemOutCtrl.UpdateItemOut).Methods("POST")
	r.HandleFunc("/item-out/delete", itemOutCtrl.DeleteItemOut).Methods("POST")
	r.HandleFunc("/item-out/export", itemOutCtrl.ExportItemOuts).Methods("POST")
	r.HandleFunc("/item-out/import", itemOutCtrl.ImportItemOuts).Methods("POST")
	// report
	r.HandleFunc("/item-value-report", reportCtrl.GetItemValueReport).Methods("GET")
	r.HandleFunc("/item-value-report/export", reportCtrl.ExportItemValueReport).Methods("POST")
	r.HandleFunc("/selling-report", reportCtrl.GetSellingReport).Methods("GET")
	r.HandleFunc("/selling-report", reportCtrl.GetSellingReport).Methods("POST")
	r.HandleFunc("/selling-report/export", reportCtrl.ExportSellingReport).Methods("POST")

	port := getPort()
	fmt.Println("Open http://localhost" + port + " to get started!")
	http.ListenAndServe(port, r)
}
