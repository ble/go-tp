package handler

import (
	"ble/tpg/room"
	"encoding/json"
	"io"
	"io/ioutil"
	. "net/http"
	"strconv"
	"time"
)

type gameHandler struct {
	room.RoomService
}

//we serve the following:
//GET <game-id>
//GET <game-id>/client
//GET <game-id>/join-client
//GET <game-id>/events
//POST <game-id>/chat
//POST <game-id>/join
//POST <game-id>/pass 
func (g *gameHandler) ServeHTTP(w ResponseWriter, r *Request) {
	parts := pathParts(r)
	if len(parts) < 1 || len(parts) > 2 {
		NotFound(w, r)
		return
	}

	gameId := parts[0]
	room, err := g.RoomService.GetRoom(gameId)
	if err != nil {
		NotFound(w, r)
		return
	}
	userId, _ := getUserId(r)
	playerId, _ := getPlayerId(r)
	bodyBytes, _ := ioutil.ReadAll(&io.LimitedReader{r.Body, 1024})
	gamePath, _ := g.RoomService.PathTo(room.GetGame())
	//split by paths:
	if len(parts) == 1 {
		if !isGet(r) {
			Error(w, "", StatusMethodNotAllowed)
			return
		}
	} else if len(parts) == 2 {

		//check the method
		switch parts[1] {
		case "join":
		case "chat":
		case "pass":
			if !isPost(r) {
				Error(w, "", StatusMethodNotAllowed)
				return
			}
		case "client":
		case "join-client":
		case "events":
			if !isGet(r) {
				Error(w, "", StatusMethodNotAllowed)
				return
			}
		}

		//actually process requests
		switch parts[1] {
		case "join":
			if pidNew, err := room.Join(userId, playerId, bodyBytes); err == nil {
				cookie := &Cookie{
					Name:     "playerId",
					Value:    pidNew,
					Path:     gamePath.String(),
					HttpOnly: true}
				w.Header().Add("Location", gamePath.String()+"/client")
				SetCookie(w, cookie)
				w.WriteHeader(StatusSeeOther)
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		case "chat":
			if err = room.Chat(userId, playerId, bodyBytes); err == nil {
				w.WriteHeader(StatusOK)
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		case "pass":
			if err = room.Pass(userId, playerId, bodyBytes); err == nil {
				w.WriteHeader(StatusOK)
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		case "client":
			if err = room.AccessClient(userId, playerId); err == nil {
				ServeFile(w, r, "./static/html/game-client.html")
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		case "join-client":
			if status, err := room.AccessJoinClient(userId, playerId); err == nil {
				ServeFile(w, r, "./static/html/join-client.html")
			} else if status == "already-allowed" {
				w.Header().Add("Location", gamePath.String()+"/client")
				w.WriteHeader(StatusSeeOther)
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		case "events":
			strLastQuery := r.URL.Query().Get("lastQuery")
			var lastQuery time.Time
			if lastQueryMillis, err := strconv.ParseInt(strLastQuery, 10, 64); err == nil {
				lastQuery = time.Unix(0, lastQueryMillis*1000)
			} else {
				lastQuery = time.Unix(0, 0)
			}
			if events, err := room.GetEvents(userId, playerId, lastQuery); err == nil {
				respBody, err := json.Marshal(events)
				if err != nil {
					Error(w, err.Error(), StatusInternalServerError)
				} else {
					w.WriteHeader(StatusOK)
					w.Write(respBody)
				}
			} else {
				Error(w, err.Error(), StatusBadRequest)
			}
		default:
			NotFound(w, r)
		}
	} else {
		NotFound(w, r)
	}
}
