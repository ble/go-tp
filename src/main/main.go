package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func dieOnError(err error, stuff ...interface{}) {
	if err != nil {
		slice := append([]interface{}{err, ": "}, stuff...)
		log.Fatal(slice...)
	}
}

func main() {
	magicRowNumber := 42
	theCount := -1

	dbFileName := "counter.sqlite"
	db, err := sql.Open("sqlite3", dbFileName)
	dieOnError(err, "opening database")
	defer db.Close()

	result, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS count (x INTEGER PRIMARY KEY, c INTEGER)")
	dieOnError(err, "creating table if necessary")

	readRow, err := db.Prepare("SELECT * FROM count WHERE x == ?")
	dieOnError(err, "preparing readRow")
	rows, err := readRow.Query(magicRowNumber)
	dieOnError(err, "querying readRow")
	isThereARow := rows.Next()

	if !isThereARow {
		insertRow, err := db.Prepare("INSERT INTO count VALUES (?, ?);")
		dieOnError(err, "preparing insertRow")

		theCount = 0
		result, err = insertRow.Exec(magicRowNumber, theCount)
		dieOnError(err, "inserting first row")

		rows, err = readRow.Query(magicRowNumber)
		dieOnError(err, "querying readRow")
		isThereARow = rows.Next()
	}

	if !isThereARow {
		dieOnError(errors.New("row doesn't exist despite being created or previously existing?"))
	}
	dieOnError(rows.Scan(&magicRowNumber, &theCount), "scanning readRow")
	rows.Close()

	updateRow, err := db.Prepare("UPDATE count SET c = ? WHERE x = ?")
	dieOnError(err, "preparing updateRow")

	theCount++
	result, err = updateRow.Exec(theCount, magicRowNumber)
	dieOnError(err, "updating the row")

	rows, err = readRow.Query(magicRowNumber)
	dieOnError(err, "querying readRow")
	isThereARow = rows.Next()

	if !isThereARow {
		dieOnError(errors.New("row doesn't exist despite being created or previously existing?"))
	}
	dieOnError(rows.Scan(&magicRowNumber, &theCount))
	rows.Close()
	log.Print("id: ", magicRowNumber, " count: ", theCount)
}
