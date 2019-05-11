package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indam-m/ss-assessment-toko_ijah/db"
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

func main() {
	db.Open()
	// Declare a new router
	r := newRouter()

	// This is where the router is useful, it allows us to declare methods that
	// this path will be valid for
	r.HandleFunc("/item", handler).Methods("GET")
	r.HandleFunc("/item/create", handler).Methods("POST")
	r.HandleFunc("/item-in", handler).Methods("GET")
	r.HandleFunc("/item-in/create", handler).Methods("POST")
	r.HandleFunc("/item-out", handler).Methods("GET")
	r.HandleFunc("/item-out/create", handler).Methods("POST")
	r.HandleFunc("/item-value-report", handler).Methods("GET")
	r.HandleFunc("/sales-report", handler).Methods("GET")

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	http.ListenAndServe(":9876", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
