package handler

import (
	. "ble/game"
	"encoding/json"
	"fmt"
	. "net/http"
	"strings"
)

type handlerJoin struct {
	GameAgent
}

func (h handlerJoin) ServeHTTP(w ResponseWriter, r *Request) {
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(StatusMethodNotAllowed)
		w.Write([]byte("only POST"))
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "" {
		w.WriteHeader(406)
		w.Write([]byte("not json"))
		return
	}

	sent := new(Event)
	err := json.NewDecoder(r.Body).Decode(sent)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("couldn't decode"))
		return
	}

	artist, err := h.GameAgent.AddArtist(sent.Name)
	fmt.Println(artist)
	fmt.Println(err)
	if err != nil {
		sent.EventType = "Error"
		sent.Error = err.Error()
		err = json.NewEncoder(w).Encode(sent)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("couldn't encode"))
			return
		}
		w.WriteHeader(200)
		return
	}
	sent.Who = artist.Id
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sent)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(StatusOK)
	fmt.Println("asdf")
}
