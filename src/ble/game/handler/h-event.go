package handler

import (
	. "ble/game"
	"encoding/json"
	. "net/http"
	"strings"
	//  "net/url"
	"strconv"
	"time"
)

type handlerEvents struct {
	GameAgent
}

type eventsResponse struct {
	QueryTime int64       `json:"queryTime"`
	Events    []GameEvent `json:"events,omitempty"`
}

func (h handlerEvents) ServeHTTP(w ResponseWriter, r *Request) {
	if strings.ToUpper(r.Method) != "GET" {
		w.WriteHeader(StatusMethodNotAllowed)
		return
	}

	var lastQuery time.Time
	strLastQuery := r.URL.Query().Get("lastQuery")
	if strLastQuery != "" {
		lastQueryNanos, err := strconv.ParseInt(strLastQuery, 10, 64)
		if err != nil {
			w.WriteHeader(StatusInternalServerError)
			return
		}
		lastQuery = time.Unix(0, lastQueryNanos)
	}
	events, time := h.GameAgent.GetGameEvents(lastQuery)
	json.NewEncoder(w).Encode(eventsResponse{time.UnixNano(), events})
}
