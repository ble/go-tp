package handler

import (
	. "ble/game"
	"encoding/json"
	. "net/http"
)

type handlerState struct {
	a GameAgent
}

func (h handlerState) ServeHTTP(w ResponseWriter, r *Request) {
	artistId, error := getExistingArtistId(h.a, r)
	if error == nil && artistId != "" {
		view, error := h.a.View()

		if error != nil {
			Error(w, error.Error(), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		bytes, error := json.Marshal(view)

		if error != nil {
			Error(w, error.Error(), 500)
			return
		}

		w.Write(bytes)
	}
}
