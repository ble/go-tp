package handler

import (
	. "ble/game"
	. "net/http"
)

type handlerClient struct {
	a GameAgent
}

func (h handlerClient) ServeHTTP(w ResponseWriter, r *Request) {
	artistId, error := getExistingArtistId(h.a, r)

	if error != nil {
		Error(w, error.Error(), 500)
		return
	}

	if artistId == "" {
		if started, _ := h.a.IsStarted(); started {
			w.WriteHeader(200)
			w.Write([]byte("Game already in progress :("))
			return
		}
	}

	ServeFile(w, r, "./static/html/game-client.html")

}
