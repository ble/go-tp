package main

import "net/url"
import "regexp"
import "path"
import "fmt"

type Regexp *regexp.Regexp


type Route interface {
  Parses(url.URL) bool
  AsURL() url.URL
}

type RouteRoot struct { } 
func (_ RouteRoot) Parses(u url.URL) bool {
  return u.Path == "/"
}
func (_ RouteRoot) AsURL() url.URL {
  return url.URL{Path: "/"}
}

type RouteRoom struct { RId string }
var roomPattern = regexp.MustCompile("^/room/([0-9a-f]*)$")
func (r *RouteRoom) Parses(u url.URL) bool {
  submatch := roomPattern.FindStringSubmatch(u.Path)
  if submatch == nil {
    return false
  }
  r.RId = submatch[1]
  return true
}
func (r *RouteRoom) AsURL() url.URL {
  return url.URL{Path: path.Join("/room", r.RId)} 
}

type RouteDrawing struct { RouteRoom; DId string }
var drawingPattern = regexp.MustCompile("^/room/([0-9a-f]*)/drawing/([0-9a-f]*)$")
func (d *RouteDrawing) Parses(u url.URL) bool {
  submatch := drawingPattern.FindStringSubmatch(u.Path)
  if submatch == nil {
    return false
  }
  d.RId = submatch[1]
  d.DId = submatch[2]
  return true
}
func (d *RouteDrawing) AsURL() url.URL {
  return url.URL{Path: path.Join("/room", d.RId, "drawing", d.DId)}
}

type RouteSequence struct { RouteRoom; SId string }
var sequencePattern = regexp.MustCompile("^/room/([0-9a-f]*)/sequence/([0-9a-f]*)$")
func (d *RouteSequence) Parses(u url.URL) bool {
  submatch := sequencePattern.FindStringSubmatch(u.Path)
  if submatch == nil {
    return false
  }
  d.RId = submatch[1]
  d.SId = submatch[2]
  return true
}
func (s *RouteSequence) AsURL() url.URL {
  return url.URL{Path: path.Join("/room", s.RId, "sequence", s.SId)}
}

func AsRoute(u url.URL) Route {
  possible := []Route{RouteRoot{}, &RouteRoom{}, &RouteDrawing{}, &RouteSequence{}}
  for _, x := range possible {
    if x.Parses(u) {
      return x
      fmt.Println("okay")
    }
  }
  return nil
}
