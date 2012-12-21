package room

import (
	"ble/tpg/model"
	"ble/tpg/persistence"
	"ble/web"
	"errors"
	"fmt"
	"net/url"
	. "time"
)

type roomService struct {
	switchboard web.Switchboard
	backend     *persistence.Backend
	rooms       map[string]Room
}

func (r *roomService) GetRoom(gameId string) (Room, error) {
	game, present := r.backend.AllGames()[gameId]
	if !present {
		return nil, errors.New("no such game")
	}
	eventsIn, eventRequestsIn := make(chan interface{}), make(chan interface{})
	room := &aRoom{r, eventsIn, eventRequestsIn, game}
	go room.processEvents()
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

func (a *aRoom) processEvents() {
	panic("not implemented")
}

func (a *aRoom) GetGame() model.Game {
	return a.game
}

func (a *aRoom) GetState() interface{} {
	return []interface{}{}
}

func (a *aRoom) Chat(uid, pid string, body []byte) error {
	return errors.New("unimplemented")
}

func (a *aRoom) Join(uid, pid string, body []byte) (string, error) {
	return "", errors.New("unimplemented")
}

func (a *aRoom) Pass(uid, pid string, body []byte) error {
	return errors.New("unimplemented")
}

func (a *aRoom) AccessClient(uid, pid string) error {
	return errors.New("unimplemented")
}

func (a *aRoom) AccessJoinClient(uid, pid string) (string, error) {
	return "", errors.New("unimplemented")
	//return "already-allowed", errors.New("can't join again")
}

func (a *aRoom) GetEvents(uid, pid string, lastQuery Time) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func typecheckRoom() Room {
	return &aRoom{}
}

func typecheckRoomService() RoomService {
	return &roomService{}
}
