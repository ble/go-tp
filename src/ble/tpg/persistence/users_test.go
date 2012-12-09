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
	err = createTables(b)
	if err != nil {
		t.Fatal(err)
	}
	u, err := b.CreateUser("the.bomb@thebomb.com", "scatman juan", "asdfquxl")
	if err != nil || u == nil {
		t.Fatal("couldn't create user")
	}
	u2, err := b.LogInUser(u.Alias(), "asdfquxll")
	if err == nil || u2 != nil {
		t.Fatal("logged in with a bad password")
	}
	u3, err := b.LogInUser(u.Alias(), "asdfquxl")
	if err != nil || u3 == nil {
		t.Fatal("failed to log in")
	}
}
