package handler

import (
	"ble/game"
	"bytes"
	"encoding/json"
	"io/ioutil"
	. "net/http"
	. "testing"
)

func hEventF(g game.GameAgent) Handler {
	return handlerEvents{g}
}

func TestJsonEvents(t *T) {
	th := NewHarness(t)
	th.start(hEventF)
	defer th.stop()

	biff, _ := th.GameAgent.AddArtist("biffhh")
	sally, _ := th.GameAgent.AddArtist("sally")
	th.GameAgent.AddArtist("biffhh")
	wizzud, _ := th.GameAgent.AddArtist("wizzud")
	th.GameAgent.Start()
	th.GameAgent.PassSequence(biff.Id)
	th.GameAgent.PassSequence(sally.Id)
	th.GameAgent.PassSequence(sally.Id)
	th.GameAgent.PassSequence(biff.Id)
	th.GameAgent.PassSequence(wizzud.Id)

	client := plainClient()
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
