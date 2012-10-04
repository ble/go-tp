package main

import (
  "net/http"
  "crypto/sha1"
  "hash"
  "time"
  "bytes"
  "encoding/base64"
  "fmt"
)


type AppState struct {
  Users map[string] bool
  Rooms map[string] *Room
}

func (a *AppState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  route := GetRoute(r, a)
  if route == nil {
    http.Error(w, "", 404)
    return
  }
  route.ServeHTTP(w, r)
  for k, v := range(a.Rooms) {
    fmt.Println(k)
    fmt.Println(v)
  }
}

func NewAppState() *AppState {
  r := new(AppState)
  r.Users = make(map[string]bool)
  r.Rooms = make(map[string]*Room)
  return r
}

func (a *AppState) CreateUserId(r *http.Request) string {
  hasher := sha1.New()
  r.Write(hasher)
  hasher.Write(timeBytes())
  id := hashToBase64(hasher)

  if a.UserExists(id) {
    return a.CreateUserId(r)
  }

  a.AssignUserId(id)
  return id
}

func (a *AppState) AssignUserId(id string) {
  a.Users[id] = true
}

func (a *AppState) UserExists(id string) bool {
  return a.Users[id]
}

func (a *AppState) UserCanCreateRoom(id string) bool {
  return a.UserExists(id)
}

func (a *AppState) CreateRoom(r *http.Request, userId string, j map[string]interface{}) string {
  hasher := sha1.New()
  r.Write(hasher)
  hasher.Write(timeBytes())
  id := hashToBase64(hasher)
  if a.RoomExists(id) {
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

func (a *AppState) RoomExists(id string) bool {
  return a.Rooms[id] != nil
}

func (a *AppState) AssignRoomId(r *Room, id string) {
  a.Rooms[id] = r
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

func timeBytes() []byte {
  nanos := time.Now().Unix()
  bytes := make([]byte, 8, 8)
  for i := 0; i < 8; i++ {
    var shift uint = uint(i << 3)
    bytes[i] = byte((nanos & (255 << shift)) >> shift)
  }
  return bytes 
}

func hashToBase64(h hash.Hash) string {
  buffer := new(bytes.Buffer)
  base64.NewEncoder(base64.URLEncoding, buffer).Write(h.Sum(nil))
  return buffer.String()

}
