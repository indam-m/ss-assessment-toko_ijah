package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open() {
	db, err := sql.Open("sqlite3", "tokoijah.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// test connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
