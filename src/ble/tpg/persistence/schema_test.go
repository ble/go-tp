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
	rs, es := createTables(backend)
	for ix := range rs {
		//r := rs[ix]
		e := es[ix]
		if e == nil {
			continue
		}
		s := tableCreationStatements[ix]
		t.Log(e, s)
	}
	t.Log(rs)
	t.Log(es)
	t.Log("HI!")
}
