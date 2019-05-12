package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

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

func getStringDate(str string) string {
	t, err := time.Parse("2006/01/02 15:04:05", str)
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	return t.Format(time.RFC3339)
}

func convertToStr(val int64) string {
	str := strconv.FormatInt(val, 10)
	return str
}

func getDateTimeStr(str string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str
	}
	return t.Format("2006/01/02 15:04:05")
}
