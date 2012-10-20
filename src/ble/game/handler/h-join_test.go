package handler

import (
	. "ble/game"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	. "testing"
)

var join handlerFactory = func(g GameAgent) http.Handler {
	return handlerJoin{g}
}

func Test_Join_Handler(t *T) {
	testMethods := func(t *T, g GameAgent, url string, c http.Client) {
		responseGet, _ := c.Get(url)
		if responseGet.StatusCode != http.StatusMethodNotAllowed {
			t.Error("non-POST method allowed")
		}

		responseHead, _ := c.Head(url)
		if responseHead.StatusCode != http.StatusMethodNotAllowed {
			t.Error("non-POST method allowed")
		}

	}

	theTest := func(t *T, g GameAgent, url string, c http.Client) {
		//add artist
		name0 := "sammy"
		postBody, _ := json.Marshal(JoinEvent(name0, ""))
		t.Log("sending(0): ", string(postBody))
		response, _ := c.Post(url, "application/json", bytes.NewReader(postBody))

		if response.StatusCode != http.StatusOK {
			t.Log(response.StatusCode, response.Status)
			t.Error("error on legit post")
			bytes, _ := ioutil.ReadAll(response.Body)
			t.Log("body: ", string(bytes))
		}
		received := new(Event)
		err := json.NewDecoder(response.Body).Decode(received)
		if err != nil {
			t.Fatal("failed to unmarshal response")
		}

		//would be nice to have convenience methods around these bodies...
		bs, _ := json.Marshal(received)
		t.Log("receiving(0): ", string(bs))
		if received.Name != name0 {
			t.Error("name doesn't match that submitted")
		}
		if received.Who == "" {
			t.Error("no id in response")
		}
		if received.EventType != "JoinGame" {
			t.Error("response is not the right event type")
		}

		//add duplicate artist
		t.Log("sending(1): ", string(postBody))
		response, _ = c.Post(url, "application/json", bytes.NewReader(postBody))
		if response.StatusCode != http.StatusOK {
			t.Error("error on legit post")
		}
		received = new(Event)
		err = json.NewDecoder(response.Body).Decode(received)
		if err != nil {
			t.Fatal("failed to unmarshal response")
		}
		bs, _ = json.Marshal(received)
		t.Log("receiving(1): ", string(bs))

		if received.EventType != "Error" {
			t.Error("duplicate name does not cause error")
		}
		if received.Error == "" {
			t.Error("no error string on duplicate name")
		}

		//add artist after starting game
		g.Start()

		name1 := "smackdab"
		postBody, _ = json.Marshal(JoinEvent(name1, ""))
		t.Log("sending(2): ", string(postBody))

		response, _ = c.Post(url, "application/json", bytes.NewReader(postBody))
		if response.StatusCode != http.StatusOK {
			t.Error("error on legit post")
		}

		received = new(Event)
		err = json.NewDecoder(response.Body).Decode(received)
		if err != nil {
			t.Fatal("failed to unmarshal response")
		}
		bs, _ = json.Marshal(received)
		t.Log("receiving(2): ", string(bs))

		if received.EventType != "Error" {
			t.Error("post-start join does not cause error")
		}
		if received.Error == "" {
			t.Error("no error string on post-start join")
		}

	}
	testInContext(t, join, testMethods)
	testInContext(t, join, theTest)
}

type handlerFactory func(GameAgent) http.Handler
type testAction func(t *T, g GameAgent, url string, c http.Client)

func testInContext(t *T, h handlerFactory, test testAction) {
	agent := NewGame()
	defer agent.Shutdown()

	server := httptest.NewServer(h(agent))
	defer server.Close()

	client := http.Client{}
	test(t, agent, server.URL, client)
}
