package persistence

import (
	DD "ble/drawing"
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

	theStack := g.Stacks()[0]
	theDrawing := theStack.AllDrawings()[0]
	if err := theDrawing.Add(DD.DefaultDrawPart); err != nil {
		t.Fatal(err)
	}
	if err := theDrawing.Add(DD.DefaultDrawPart); err != nil {
		t.Fatal(err)
	}
	if err := theDrawing.Complete(); err != nil {
		t.Fatal(err)
	}
	if err := g.PassStack(p); err != nil {
		t.Fatal(err)
	}
	theStack.AddDrawing(p)
	t.Logf("%#v", theStack)
	t.Logf("%#v", theDrawing)
	t.Logf("%#v", theStack.AllDrawings()[1])
}
