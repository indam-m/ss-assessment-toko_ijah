package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
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

func convertDateForSQL(str string) string {
	re := regexp.MustCompile("(\\d{4})[/-](\\d{1,2})[/-](\\d{1,2}) (\\d{1,2}):(\\d{1,2})")
	var locStr string
	match := re.FindStringSubmatch(str)
	for i, val := range match {
		if i > 0 {
			if len(val) == 1 {
				locStr += "0"
			}
			locStr += val
			if i < 3 {
				locStr += "/"
			} else if i > 3 {
				locStr += ":"
			} else {
				locStr += " "
			}
		}
	}
	// adding second
	re = regexp.MustCompile("(\\d{1,2}):(\\d{1,2}):(\\d{1,2})")
	match = re.FindStringSubmatch(str)
	if len(match) > 0 {
		sec := match[3]
		if len(sec) == 1 {
			locStr += "0"
		}
		locStr += sec
	} else {
		locStr += "00"
	}

	t, err := time.Parse("2006/01/02 15:04:05", locStr)
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	return t.Format(time.RFC3339)
}

func convertToStr(val int64) string {
	str := strconv.FormatInt(val, 10)
	return str
}

func convertToInt(str string) int64 {
	re := regexp.MustCompile(",")
	str = re.ReplaceAllString(str, "")
	if len(str) > 2 && str[0:2] == "Rp" {
		str = str[2:]
	}
	res, _ := strconv.ParseInt(str, 10, 64)
	return res
}

func convertToUITime(str string) string {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str
	}
	return t.Format("2006/01/02 15:04:05")
}

func execImport(sqlStr string, vals []interface{}, w http.ResponseWriter) error {
	// prepare the statement
	stmt, err := database.Prepare(sqlStr)
	if err != nil {
		fmt.Fprintln(w, "Prepare query error")
		fmt.Fprintf(w, err.Error())
		return err
	}
	// execute the statement
	_, err = stmt.Exec(vals...)
	if err != nil {
		fmt.Fprintln(w, "Execute query error")
		fmt.Fprintf(w, err.Error())
		return err
	}
	return err
}
