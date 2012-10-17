package handler

import (
	. "ble/game"
	. "net/http"
)

type roomHandler struct {
	agent GameAgent
}

//room routes are:
//   -prefix-/client
//   -prefix-/state
//   -prefix-/notification
//   -prefix-/drawing/#id#
func (h roomHandler) ServeHTTP(w ResponseWriter, r *Request) {

}

func getExistingArtistId(g GameAgent, r *Request) (string, error) {
	cookie, error := r.Cookie("artistId")
	if error != nil {
		return "", error
	}
	present, error := g.HasArtist(cookie.Value)
	if error == nil && present {
		return cookie.Value, nil
	}
	return "", error
}
