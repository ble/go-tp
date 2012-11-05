package handler

import (
	"ble/game"
	. "ble/testing/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	. "net/http"
	. "testing"
)

func TestJsonEvents(t *T) {
	g := game.NewGame()
	th := NewHarness(t, FromHandler(handlerEvents{g}))
	defer th.Stop()

	biff, _ := g.AddArtist("biffhh")
	sally, _ := g.AddArtist("sally")
	g.AddArtist("biffhh")
	wizzud, _ := g.AddArtist("wizzud")
	g.Start()
	g.PassSequence(biff.Id)
	g.PassSequence(sally.Id)
	g.PassSequence(sally.Id)
	g.PassSequence(biff.Id)
	g.PassSequence(wizzud.Id)

	client := PlainClient()
	url := th.URL.String()
	req, _ := NewRequest("GET", url, nil)
	resp, _ := client.Do(req)
	oJson := &eventsResponse{0, make([]game.GameEvent, 0, 0)}
	b, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(b))
	_ = json.NewDecoder(bytes.NewReader(b)).Decode(oJson)
	for _, v := range oJson.Events {
		bs, _ := json.Marshal(v)
		t.Log(string(bs))
	}
}
