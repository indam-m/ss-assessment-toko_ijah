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

	// declaring routers
	r.HandleFunc("/item-amount", itemAmountCtrl.GetItemAmounts).Methods("GET")
	r.HandleFunc("/item-amount/{sku}", func(w http.ResponseWriter, r *http.Request) {
		itemAmountCtrl.GetItemAmount(w, r, getSKU(r))
	}).Methods("GET")
	r.HandleFunc("/item-amount/create", itemAmountCtrl.CreateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/update", itemAmountCtrl.UpdateItemAmount).Methods("POST")
	r.HandleFunc("/item-amount/delete", itemAmountCtrl.DeleteItemAmount).Methods("POST")

	r.HandleFunc("/item-in", itemInCtrl.GetItemIns).Methods("GET")
	r.HandleFunc("/item-in/create", itemInCtrl.CreateItemIn).Methods("POST")
	r.HandleFunc("/item-in/update", itemInCtrl.UpdateItemIn).Methods("POST")
	r.HandleFunc("/item-in/delete", itemInCtrl.DeleteItemIn).Methods("POST")
	r.HandleFunc("/item-out", handler).Methods("GET")
	r.HandleFunc("/item-out/create", handler).Methods("POST")
	r.HandleFunc("/item-value-report", handler).Methods("GET")
	r.HandleFunc("/sales-report", handler).Methods("GET")

	http.ListenAndServe(":9876", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
