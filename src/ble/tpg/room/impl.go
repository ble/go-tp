package room

import (
	"ble/tpg/model"
	"ble/tpg/persistence"
	"ble/web"
	"errors"
	"fmt"
	"net/url"
)

type roomService struct {
	switchboard web.Switchboard
	backend     *persistence.Backend
	rooms       map[string]Room
}

func NewRoomService(s web.Switchboard, b *persistence.Backend) RoomService {
	return &roomService{
		switchboard: s,
		backend:     b,
		rooms:       make(map[string]Room)}
}

func (r *roomService) GetRoom(gameId string) (Room, error) {
	game, present := r.backend.AllGames()[gameId]
	if !present {
		return nil, errors.New("no such game")
	}
	room, present := r.rooms[gameId]
	if present {
		return room, nil
	}
	eventsIn, eventRequestsIn := make(chan interface{}), make(chan interface{})
	newRoom := &aRoom{r, eventsIn, eventRequestsIn, game}
	go newRoom.processEvents()
	r.rooms[gameId] = newRoom
	return room, nil
}

func (r *roomService) PathTo(obj interface{}) (*url.URL, error) {
	url := r.switchboard.URLOf(obj)
	if url == nil {
		return nil, fmt.Errorf("don't know path for object %#v", obj)
	}
	return url, nil
}

type aRoom struct {
	*roomService
	events        chan interface{}
	eventRequests chan interface{}
	game          model.Game
}

func (a *aRoom) GetGame() model.Game {
	return a.game
}

func (a *aRoom) GetState() interface{} {
	return gameJson{a.game, a.roomService.switchboard}
}

func (a *aRoom) Chat(uid, pid string, body []byte) error {
	action, errC := aschat(body)
	if a.game.PlayerForId(pid) != nil &&
		errC == nil {
		a.events <- chat(pid, action.Content)
		return nil
	}
	return errors.New("")

}

func (a *aRoom) Join(uid, pid string, body []byte) (string, error) {
	user, err := a.roomService.backend.GetUserById(uid)
	if err != nil {
		return "", errors.New("not logged in")
	}
	action, err := asjoingame(body)
	if err != nil {
		return "", errors.New("bad request")
	}
	player, err := a.game.JoinGame(user, action.Name)
	if err != nil {
		return "", err
	}
	a.events <- joingame(player.Pid(), action.Name)
	return player.Pid(), nil
}

func (a *aRoom) Pass(uid, pid string, body []byte) error {
	player := a.game.PlayerForId(pid)
	if player == nil {
		return errors.New("no such player")
	}
	_, err := aspassstack(body)
	if err != nil {
		return err
	}
	stack, err := a.game.PassStack(player)
	if err != nil {
		return err
	}

	if stack.TopDrawing().Player() != player {
		nextPlayer := a.game.NextPlayer(player)
		a.events <- passstack(player.Pid(), nextPlayer.Pid(), stack.Sid())
	} else {
		a.events <- passstack(player.Pid(), "", stack.Sid())
	}
	return nil
}

func (a *aRoom) Start(uid, pid string, body []byte) error {
	player := a.game.PlayerForId(pid)
	if player == nil {
		return errors.New("no such player")
	}
	_, err := asstartgame(body)
	if err != nil {
		return err
	}
	err = a.game.Start()
	if err != nil {
		return err
	}
	a.events <- startgame(pid)
	return nil
}

func (a *aRoom) AccessClient(uid, pid string) error {
	if a.game.PlayerForId(pid) != nil {
		return nil
	}
	return errors.New("no such player")
}

func (a *aRoom) AccessJoinClient(uid, pid string) (string, error) {
	if _, err := a.backend.GetUserById(uid); err != nil {
		return "", errors.New("not logged in")
	}
	if a.game.PlayerForId(pid) != nil {
		return "already-allowed", errors.New("can't join again")
	}
	if !a.game.IsStarted() {
		return "", nil
	}
	return "", errors.New("game already started")
}

func typecheckRoom() Room {
	return &aRoom{}
}

func typecheckRoomService() RoomService {
	return &roomService{}
}
