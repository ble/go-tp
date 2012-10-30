package handler

import (
	. "ble/game"
	"encoding/json"
	. "net/http"
	"net/url"
	"path"
	"strings"
)

type handlerJoin struct {
	GameAgent
}

func (h handlerJoin) ServeHTTP(w ResponseWriter, r *Request) {
	//reject non-post
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(StatusMethodNotAllowed)
		w.Write([]byte("only POST"))
		return
	}

	//reject overly-long post
	if r.ContentLength >= 1024 {
		w.WriteHeader(StatusRequestEntityTooLarge)
		return
	}

	//accept only json
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "" {
		w.WriteHeader(StatusUnsupportedMediaType)
		w.Write([]byte("not json"))
		return
	}

	//try to interpret json
	sent := new(Event)
	err := json.NewDecoder(r.Body).Decode(sent)
	if err != nil {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte("couldn't decode"))
		return
	}

	//reject if cookie shows that one has already joined
	existingId, err := getExistingArtistId(h.GameAgent, r)
	if existingId != "" {
		w.WriteHeader(StatusOK)
		errorResponse := new(Event)
		errorResponse.EventType = "Error"
		errorResponse.Error = "you've already joined this game"
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	//call into game
	artist, err := h.GameAgent.AddArtist(sent.Name)
	if err != nil {

		//try to send descriptive error
		sent.EventType = "Error"
		sent.Error = err.Error()
		jsonBytes, err := json.Marshal(sent)

		//somehow fail to encode JSON
		if err != nil {
			w.WriteHeader(StatusInternalServerError)
			w.Write([]byte("couldn't encode"))
			return
		}
		w.WriteHeader(StatusOK)
		w.Write(jsonBytes)
		return
	}

	//call into game succeeded;

	//make the cookie
	//we assume that request path is like /path/to/room/join
	//and that in such a case, the path of all room resources is
	//path/to/room.
	cookiePath, err := url.Parse(path.Dir(r.URL.Path))
	aIdCookie := &Cookie{
		Name:     "artistId",
		Value:    artist.Id,
		Path:     cookiePath.String(),
		HttpOnly: true}
	SetCookie(w, aIdCookie)
	w.WriteHeader(StatusOK)

	//prepare the response body
	sent.Who = artist.Id
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sent)
}
