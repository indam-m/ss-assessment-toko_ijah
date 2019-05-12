package controller

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // library to open sqlite3 database
)

var (
	database *sql.DB
	dbErr    error
)

// Open opens the tokoijah sqlite database
func Open() {
	database, dbErr = sql.Open("sqlite3", "./tokoijah.db")
	if dbErr != nil {
		panic(dbErr)
	}
	// defer database.Close()
	// test connection
	dbErr = database.Ping()
	if dbErr != nil {
		panic(dbErr)
	}
}

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
