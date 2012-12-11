package persistence

import (
	_ "github.com/ble/go-sqlite3"
	"os"
	. "testing"
)

func TestCreateTables(t *T) {
	backend, err := NewBackend("testdb")
	defer os.Remove("testdb")
	defer backend.conn.Close()
	t.Log(err)
	err = createTables(backend)
	if err != nil {
		t.Fatal(err)
	}
	err = backend.prepAllStatements()
	if err != nil {
		t.Fatal(err)
	}
}
