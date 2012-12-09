package persistence

import (
	"os"
	. "testing"
)

func TestCreateUser(t *T) {
	b, err := NewBackend("testdb")
	b.RegisterLogger(t)
	t.Log(err)
	defer os.Remove("testdb")
	rs, es := createTables(b)
	t.Log(rs, es)
	u, err := b.CreateUser("the.bomb@thebomb.com", "scatman juan", "asdfquxl")
	t.Log(err, u)
	u2, err := b.LogInUser(u.Alias(), "asdfquxll")
	t.Log(err, u2)
	u3, err := b.LogInUser(u.Alias(), "asdfquxl")
	t.Log(err, u3)
}
