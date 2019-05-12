package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indam-m/ss-assessment-toko_ijah/controller"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	assets := "/assets/"
	staticFileDirectory := http.Dir(assets)
	staticFileHandler := http.StripPrefix(assets, http.FileServer(staticFileDirectory))
	r.PathPrefix(assets).Handler(staticFileHandler).Methods("GET")

	return r
}

func getSKU(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["sku"]
}

func main() {
	controller.Open()
	// Declare a new router
	r := newRouter()

	itemAmountCtrl := &controller.ItemAmount{}
	itemInCtrl := &controller.ItemIn{}
	itemOutCtrl := &controller.ItemOut{}
	reportCtrl := &controller.Report{}

	// declaring routers
	// item amount
	r.HandleFunc("/item-amount", itemAmountCtrl.GetItemAmounts).Methods("GET")
	r.HandleFunc("/item-amount/{sku}", func(w http.ResponseWriter, r *http.Request) {
		itemAmountCtrl.GetItemAmount(w, r, getSKU(r))
	}).Methods("GET")
	r.HandleFunc("/item-amount/create", itemAmountCtrl.CreateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/update", itemAmountCtrl.UpdateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/delete", itemAmountCtrl.DeleteItemAmount).Methods("POST")
	// item in
	r.HandleFunc("/item-in", itemInCtrl.GetItemIns).Methods("GET")
	r.HandleFunc("/item-in/create", itemInCtrl.CreateItemIn).Methods("POST")
	r.HandleFunc("/item-in/update", itemInCtrl.UpdateItemIn).Methods("POST")
	r.HandleFunc("/item-in/delete", itemInCtrl.DeleteItemIn).Methods("POST")
	// item out
	r.HandleFunc("/item-out", itemOutCtrl.GetItemOuts).Methods("GET")
	r.HandleFunc("/item-out/create", itemOutCtrl.CreateItemOut).Methods("POST")
	r.HandleFunc("/item-out/update", itemOutCtrl.UpdateItemOut).Methods("POST")
	r.HandleFunc("/item-out/delete", itemOutCtrl.DeleteItemOut).Methods("POST")
	// report
	r.HandleFunc("/item-value-report", reportCtrl.GetItemValueReport).Methods("POST")
	r.HandleFunc("/selling-report", reportCtrl.GetSellingReport).Methods("POST")

	http.ListenAndServe(":9876", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
