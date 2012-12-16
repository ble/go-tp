package persistence

import (
	_ "github.com/ble/go-sqlite3"
	"os"
	. "testing"
)

func TestCreateTables(t *T) {
	backend, err := NewBackend("testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testdb")
	defer backend.conn.Close()
}
