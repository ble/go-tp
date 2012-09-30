package main

import (
  "log"
  "net/http"
  "fmt"
)

type Nexus struct {
  rooms map[string] *Room
}

func (n *Nexus) addRoom(roomName string) {
  newRoom := &Room{"asdf", roomName, make(map[string]int)}
  n.rooms[newRoom.Id] = newRoom
}

func (n *Nexus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  roomId := GetRoom(r.URL)
  uk, ukError := GetUserKey(r)
  if ukError != nil {
    
  }
  fmt.Println(roomId)
  fmt.Println(uk)
}

func main() {
  drawing := NewDrawing(640, 480, "nathaniel hawthorne")
  nexus := &Nexus{make(map[string]*Room)}
  http.Handle("/nat", &drawing)
  http.Handle("/", nexus)

  log.Fatal(http.ListenAndServe(":8080", nil))
}
