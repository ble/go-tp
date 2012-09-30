package main

import (
  "errors"
)

type Room struct {
  Id, Name string
  Artists map[string]int
}

type RoomView struct {
  *Room
  ArtistId string
}

func (r *Room) GetView(u UserKey) (*RoomView, error) {
  if r.Id != u.Room {
    return nil, errors.New("wrong room")
  }
  if r.Artists[u.ArtistId] != 0 {
    return nil, errors.New("not present in room")
  }
  return &RoomView{r, u.ArtistId}, nil
}
