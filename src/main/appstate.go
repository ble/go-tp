package main

import (
	. "ble/hash"
	"fmt"
	"net/http"
	"time"
)

type AppHandle struct {
	connection chan<- interface{}
}

type appState struct {
	users map[string]bool
	rooms map[string]*Room
}

type createUser struct {
	success
	id string
}
type userExists struct {
	success
	id string
}
type getRooms struct {
	success
	rooms map[string]*Room
}
type existsRoom struct {
	success
	id string
}
type createRoom struct {
	success
	id   string
	room *Room
}

func (a AppHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := GetRoute(r, a)
	if route == nil {
		http.Error(w, "", 404)
		return
	}
	route.ServeHTTP(w, r)
	rooms, _ := a.getRooms()
	for k, v := range rooms {
		fmt.Println(k)
		fmt.Println(v)
	}
}

func NewAppState() (AppHandle, appState) {
	state := appState{make(map[string]bool), make(map[string]*Room)}
	ch := make(chan interface{})
	handle := AppHandle{ch}
	go state.serve(ch)
	return handle, state
}

func (a appState) serve(connection <-chan interface{}) {
	for v := range connection {
		fmt.Printf("%#v", v)
		switch vv := v.(type) {
		case createUser:
			a.users[vv.id] = true
			vv.success <- a.users[vv.id]
		case userExists:
			fmt.Printf("User %s exists: %s", vv.id, a.users[vv.id])
			vv.success <- a.users[vv.id]
		case getRooms:
			vv.rooms = make(map[string]*Room)
			for k, v := range a.rooms {
				vv.rooms[k] = v
			}
			vv.success <- true
		case existsRoom:
			vv.success <- (a.rooms[vv.id] != nil)
		case createRoom:
			a.rooms[vv.id] = vv.room
			vv.success <- true
		default:
			fmt.Println("unknown type")
		}
	}
	fmt.Println("done")
}

func (a AppHandle) CreateUserId(r *http.Request) string {
	hasher := NewHashEasy()
	hasher.Nonce()
	r.Write(hasher)
	id := hasher.String()

	if exists, _ := a.UserExists(id); exists {
		fmt.Println("hashcol")
		return a.CreateUserId(r)
	}

	a.AssignUserId(id)
	return id
}

func (a AppHandle) AssignUserId(id string) error {
	req := createUser{make(success), id}
	a.connection <- req
	_, error := req.succeededIn(1 * time.Second)
	return error
}

func (a AppHandle) UserExists(id string) (bool, error) {
	req := userExists{make(success), id}
	a.connection <- req
	ok, error := req.succeededIn(1 * time.Second)
	return ok, error
}

func (a AppHandle) getRooms() (map[string]*Room, error) {
	req := getRooms{make(success), make(map[string]*Room)}
	a.connection <- req
	_, error := req.succeededIn(1 * time.Second)
	return req.rooms, error
}

func (a AppHandle) UserCanCreateRoom(id string) (bool, error) {
	return a.UserExists(id)
}

func (a AppHandle) CreateRoom(r *http.Request, userId string, j map[string]interface{}) string {
	hasher := NewHashEasy()
	r.Write(hasher)
	hasher.Nonce()
	id := hasher.String()
	if ex, _ := a.RoomExists(id); ex {
		return a.CreateRoom(r, userId, j)
	}
	newRoom := Room{}
	switch j["name"].(type) {
	case nil:
		newRoom.Name = "Nameless room"
	case string:
		newRoom.Name = j["name"].(string)
	default:
		newRoom.Name = "Mystery room"
	}
	newRoom.UserCreator = userId
	a.AssignRoomId(&newRoom, id)
	return id
}

func (a AppHandle) RoomExists(id string) (bool, error) {
	req := existsRoom{make(success), id}
	a.connection <- req
	ok, error := req.succeededIn(1 * time.Second)
	return ok, error
}

func (a AppHandle) AssignRoomId(r *Room, id string) error {
	req := createRoom{make(success), id, r}
	a.connection <- req
	_, error := req.succeededIn(1 * time.Second)
	return error
}

type UserInfo struct {
	UserId, ArtistId, RoomId string
}

func (info *UserInfo) ExtractFrom(request *http.Request) {
	userId, error := request.Cookie("user")
	if error == nil {
		info.UserId = userId.Value
	}
	roomId, error := request.Cookie("room")
	if error == nil {
		info.RoomId = roomId.Value
	}
	artistId, error := request.Cookie("artist")
	if error == nil {
		info.ArtistId = artistId.Value
	}

}
