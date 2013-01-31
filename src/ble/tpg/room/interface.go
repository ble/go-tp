package room

import (
	"ble/tpg/model"
	"ble/web"
	"io"
	"net/url"
	"time"
)

type RoomService interface {
	GetRoom(string) (Room, error)
	GetStackAndRoom(string) (model.Stack, Room, error)
	GetDrawingAndRoom(string) (model.Drawing, Room, error)
	PathTo(interface{}) (*url.URL, error)
	GetSwitchboard() web.Switchboard
}

type Room interface {
	GetGame() model.Game
	GetState(pid string) interface{}
	Chat(uid, pid string, body []byte) error
	Join(uid, pid string, body []byte) (string, error)
	Pass(uid, pid string, body []byte) error
	Start(uid, pid string, body []byte) error
	Draw(uid, pid string, d model.Drawing, body io.Reader) error
	AccessClient(uid, pid string) error
	AccessJoinClient(uid, pid string) (string, error)
	GetEvents(uid, pid string, lastQueryTime time.Time) (interface{}, error)
}
