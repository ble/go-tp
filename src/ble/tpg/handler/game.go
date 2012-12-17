package handler

import (
	. "ble/tpg/model"
	"ble/tpg/persistence"
	"ble/web"
	"encoding/json"
	. "net/http"
)

type gameHandler struct {
	*persistence.Backend
	switchboard web.Switchboard
}

//we serve the following:
//GET <game-id>
//GET <game-id>/client
//GET <game-id>/join-client
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
	game, present := g.AllGames()[gameId]
	if !present {
		NotFound(w, r)
	}

	//technically this is kinda backwards...
	//should determine whether something exists from the path alone,
	//then look at the method to decide if it is allowed...
	if isPost(r) {
		if len(parts) != 2 {
			NotFound(w, r)
			return
		}
		switch parts[1] {
		case "join":
			g.hJoin(game, w, r)
		case "chat":
			g.hChat(game, w, r)
		case "pass":
			g.hPass(game, w, r)
		default:
			NotFound(w, r)
		}
	} else if isGet(r) {
		if len(parts) == 1 {
			g.hGetState(game, w, r)
		} else {
			switch parts[1] {
			case "client":
				g.hClient(game, w, r)
			case "join-client":
				g.hJoinClient(game, w, r)
			}
		}
	} else {
		Error(w, "unknown method", StatusMethodNotAllowed)
	}
}

func (g *gameHandler) hClient(game Game, w ResponseWriter, r *Request) {
}
func (g *gameHandler) hGetState(game Game, w ResponseWriter, r *Request) {
}
func (g *gameHandler) hJoinClient(game Game, w ResponseWriter, r *Request) {
}
func (g *gameHandler) hPass(game Game, w ResponseWriter, r *Request) {
	playerId, err := getPlayerId(r)
	player := game.PlayerForId(playerId)
	if err != nil || player == nil {
		Error(w, "You have not joined this game", StatusBadRequest)
		return
	}

	if r.ContentLength >= 1024 {
		Error(w, "", StatusRequestEntityTooLarge)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		Error(w, "", StatusUnsupportedMediaType)
		return
	}

	var action Action
	err = json.NewDecoder(r.Body).Decode(&action)
	if err != nil || action.ActionType != "passStack" {
		Error(w, "Bad action", StatusBadRequest)
		return
	}

	stack, err = game.PassStack(player)

}
func (g *gameHandler) hChat(game Game, w ResponseWriter, r *Request) {
	playerId, err := getPlayerId(r)
	player := game.PlayerForId(playerId)
	if err != nil || player == nil {
		Error(w, "You have not joined this game", StatusBadRequest)
		return
	}

	if r.ContentLength >= 1024 {
		Error(w, "", StatusRequestEntityTooLarge)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		Error(w, "", StatusUnsupportedMediaType)
		return
	}

	var action Action
	err = json.NewDecoder(r.Body).Decode(&action)
	if err != nil || action.ActionType != "chat" {
		Error(w, "Bad action", StatusBadRequest)
		return
	}

	panic("chat not implemented")

}
func (g *gameHandler) hJoin(game Game, w ResponseWriter, r *Request) {
	var err error
	if _, err := getPlayerId(r); err == nil {
		Error(w, "You've already joined this game", StatusBadRequest)
		return
	}

	userId, err := getUserId(r)
	if err != nil {
		Error(w, "You are not logged in", StatusBadRequest)
	}

	if r.ContentLength >= 1024 {
		Error(w, "", StatusRequestEntityTooLarge)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		Error(w, "", StatusUnsupportedMediaType)
		return
	}

	var action Action
	err = json.NewDecoder(r.Body).Decode(&action)
	if err != nil || action.ActionType != "join" {
		Error(w, "Bad action", StatusBadRequest)
		return
	}

	user, err := g.GetUserById(userId)
	if err != nil {
		//since we got a bad userId from a cookie, we'll erase that cookie
		eraseCookie := &Cookie{
			Name:     "userId",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1} //negative maximum age in seconds indicates "delete now"
		SetCookie(w, eraseCookie)
		Error(w, "No such user", StatusBadRequest)
		return
	}

	player, err := game.JoinGame(user, action.Name)
	if err != nil {
		Error(w, err.Error(), StatusBadRequest)
		return
	}

	playerCookie := &Cookie{
		Name:     "playerId",
		Value:    player.Pid(),
		Path:     g.switchboard.URLOf(game).Path,
		HttpOnly: true}
	SetCookie(w, playerCookie)
	w.WriteHeader(StatusOK)
	return
}
