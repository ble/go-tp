package main

import (
  "net/http"
  "time"
  "encoding/json"
  "fmt"
  "image"
)

type Drawing struct {
  Creation, Finish *time.Time
  ArtistId string
  Parts []DrawPart
  Area image.Rectangle
}

func NewDrawing(width, height int, artistId string) Drawing {
  t := time.Now()
  return Drawing {
    Creation: &t,
    ArtistId: artistId,
    Parts: make([]DrawPart, 0, 100),
    Area: image.Rect(0, 0, width, height),
  }
}



func (d *Drawing) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  uk, _ := GetUserKey(r)

  switch r.Method {
    case "GET":
      uk = UserKey{"a", "b"}
      uk.SetCookieHere(w, r)
      encoder := json.NewEncoder(w)
      error := encoder.Encode(d)
      if error != nil {
        fmt.Println(error.Error())
      } 
    case "POST":
      decoder := json.NewDecoder(r.Body)
      var p DrawPart
      error := decoder.Decode(&p)
      if error != nil {
        http.Error(w, "", 415)
        fmt.Println(error.Error())
      } else {
        d.append(p)
        http.StatusText(200)
      }
    default:
      http.Error(w, "", 405)
  }
}

func (d *Drawing) append(p DrawPart) {
  d.Parts = append(d.Parts, p)
}

type DrawPart struct {
  Tag string
  When int64
  Who string
  What interface{}
}


