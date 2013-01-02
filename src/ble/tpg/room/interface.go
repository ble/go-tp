package room

import (
	"ble/tpg/model"
	"net/url"
	"time"
)

type RoomService interface {
	GetRoom(string) (Room, error)
	PathTo(interface{}) (*url.URL, error)
}

type Room interface {
	GetGame() model.Game
	GetState() interface{}
	Chat(uid, pid string, body []byte) error
	Join(uid, pid string, body []byte) (string, error)
	Pass(uid, pid string, body []byte) error
	Start(uid, pid string, body []byte) error
	AccessClient(uid, pid string) error
	AccessJoinClient(uid, pid string) (string, error)
	GetEvents(uid, pid string, lastQueryTime time.Time) (interface{}, error)
}