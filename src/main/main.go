package main

import (
	. "ble/game"
	"encoding/json"
	"fmt"
)

func Marshal(o interface{}) ([]byte, error) {
	return json.Marshal(o)
	//	return json.MarshalIndent(o, "", "\t")
}

func main() {
	g := NewGame()
	v, _ := g.View()
	gJson, _ := Marshal(v)
	fmt.Println(string(gJson))

	a1id, _ := g.AddArtist("asdf")
	a2id, _ := g.AddArtist("w0pak")
	a3id, _ := g.AddArtist("without_spaces")
	v, _ = g.View()
	gJson, _ = Marshal(v)
	fmt.Println(string(gJson))

	g.Start()
	v, _ = g.View()
	gJson, _ = Marshal(v)
	fmt.Println(string(gJson))
	g.PassSequence(a1id)
	g.PassSequence(a2id)
	g.PassSequence(a2id)
	g.PassSequence(a3id)
	g.PassSequence(a3id)
	g.PassSequence(a3id)
	g.PassSequence(a1id)
	g.PassSequence(a1id)
	g.PassSequence(a2id)
	v, _ = g.View()
	gJson, _ = Marshal(v)
	fmt.Println(string(gJson))

}
