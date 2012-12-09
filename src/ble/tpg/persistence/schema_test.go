package persistence

import (
	"database/sql"
	_ "github.com/ble/go-sqlite3"
	"os"
	. "testing"
)

func TestCreateTables(t *T) {
	db, err := sql.Open("sqlite3", "testdb")
	defer os.Remove("testdb")
	defer db.Close()
	t.Log(err)
	rs, es := createTables(db)
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
