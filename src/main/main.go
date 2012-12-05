package main

import (
	"code.google.com/p/gosqlite/sqlite"
	"errors"
	"log"
)

func dieOnError(err error, stuff ...interface{}) {
	if err != nil {
		slice := append([]interface{}{err, ": "}, stuff...)
		log.Fatal(slice...)
	}
}

func main() {
	dbFileName := "counter.sqlite"
	conn, err := sqlite.Open(dbFileName)
	defer conn.Close()
	dieOnError(err, "opening database")

	//Make statement creating table
	createTableIfNecessary, err := conn.Prepare("CREATE TABLE IF NOT EXISTS count (x INTEGER PRIMARY KEY, c INTEGER)")
	dieOnError(err, "preparing createTableIfNecessary")
	//Clean it up when done
	defer createTableIfNecessary.Finalize()
	//Actually run it
	createTableIfNecessary.Next()

	magicRowNumber := 42
	readRow, err := conn.Prepare("SELECT * FROM count WHERE x == ?")
	dieOnError(err, "preparing readRow")
	defer readRow.Finalize()
	dieOnError(readRow.Exec(magicRowNumber), "executing readRow")
	isThereARow := readRow.Next()

	theCount := -1
	if !isThereARow {
		writeFirstRow, err := conn.Prepare("INSERT INTO count VALUES (?, ?);")
		dieOnError(err, "preparing writeFirstRow")
		defer writeFirstRow.Finalize()
		dieOnError(writeFirstRow.Exec(magicRowNumber, 0), "executing writeFirstRow")
		writeFirstRow.Next()
		dieOnError(readRow.Exec(magicRowNumber), "executing readRow")
		isThereARow = readRow.Next()
	}

	if !isThereARow {
		dieOnError(errors.New("row doesn't exist despite being created or previously existing?"))
	}
	dieOnError(readRow.Scan(&magicRowNumber, &theCount), "scanning readRow")
	theCount++

	updateRow, err := conn.Prepare("UPDATE count SET c = ? WHERE x = ?")
	dieOnError(err, "updating row")
	defer updateRow.Finalize()
	dieOnError(updateRow.Exec(theCount, magicRowNumber))
	updateRow.Next()

	dieOnError(readRow.Exec(magicRowNumber))
	isThereARow = readRow.Next()
	if !isThereARow {
		dieOnError(errors.New("row doesn't exist despite being created or previously existing?"))
	}
	dieOnError(readRow.Scan(&magicRowNumber, &theCount))
	log.Print("id: ", magicRowNumber, " count: ", theCount)
}
