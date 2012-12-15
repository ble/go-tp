package persistence

import (
	"ble/tpg/model"
	"os"
	. "testing"
)

func TestCreateGame(t *T) {
	backend, err := NewBackend("testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testdb")

	backend.RegisterLogger(t)
	if err = backend.createTables(); err != nil {
		t.Fatal(err)
	}

	u, err := backend.CreateUser("the.bomb@thebomb.com", "scatman juan", "asdfff")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", u)

	gamesSvc := &games{backend.gamesBackend, make(map[string]model.Game)}
	g, err := gamesSvc.CreateGame("grapnal vs. dognel")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", g)

	p, err := g.JoinGame(u, "amazing peenee surprise")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", p)

	err = g.Start()
	if err != nil {
		t.Fatal(err)
	}

	for _, stack := range g.Stacks() {
		t.Logf("%#v", stack)
		for _, drawing := range stack.AllDrawings() {
			t.Logf("%#v", drawing)
		}
	}
}
