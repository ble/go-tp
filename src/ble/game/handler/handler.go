package handler

import (
	. "ble/game"
	. "net/http"
)

type roomHandler struct {
	agent GameAgent
	mux   *ServeMux
}

func strippingPrefix(s *ServeMux, p string, h Handler) {
	s.Handle(p, StripPrefix(p, h))
}

func NewRoomHandler(agent GameAgent) Handler {
	h := roomHandler{agent, NewServeMux()}
	strippingPrefix(h.mux, "/joinClient", handlerJoinClient{})
	strippingPrefix(h.mux, "/join", handlerJoin{agent})
	strippingPrefix(h.mux, "/client", handlerClient{agent})
	strippingPrefix(h.mux, "/state", handlerState{agent})
	strippingPrefix(h.mux, "/events", handlerEvents{agent})
	//  strippingPrefix(h.mux, "/drawing/", /*...*/)
	//  strippingPrefix(h.mux, "/stack/", /*...*/)
	return h
}

//room routes are:
//   -prefix-/client
//   -prefix-/state
//   -prefix-/notification
//   -prefix-/drawing/#id#
func (h roomHandler) ServeHTTP(w ResponseWriter, r *Request) {
	h.mux.ServeHTTP(w, r)
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
