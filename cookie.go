package main

import (
  "net/http"
  "net/url"
  "path"
)

type UserKey struct {
  Room, ArtistId string
}

func GetRoom(u *url.URL) string {
  uPath := u.Path
  for path.Base(uPath) != "/" && path.Base(uPath) != "." {
    uPath = path.Base(uPath)
  }
  return uPath
}

func GetUserKey(r *http.Request) (UserKey, error) {
  roomCookie, error := r.Cookie("room") 
  if error != nil {
    return UserKey{}, error
  } 
  artistCookie, error := r.Cookie("artist")
  if error != nil {
    return UserKey{}, error
  }
  
  return UserKey{roomCookie.Value, artistCookie.Value}, nil
}

func (u UserKey) SetCookieHere(w http.ResponseWriter, r *http.Request) {
  u.SetCookie(*r.URL, w, r)
}

func (u UserKey) SetCookie(path url.URL, w http.ResponseWriter, r *http.Request) {
  room := http.Cookie{Name: "room", Value: u.Room, Path: path.Path, HttpOnly: true}
  artist := http.Cookie{Name: "artist", Value: u.ArtistId, Path: path.Path, HttpOnly: true}
  w.Header().Add("Set-Cookie", room.String())
  w.Header().Add("Set-Cookie", artist.String())
}


