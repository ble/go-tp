package switchboard

import (
	"ble/tpg/ephemeral"
	"ble/tpg/model"
	"net/http"
)

type gameMapping struct{ http.Handler }
type stackMapping struct{ http.Handler }
type drawingMapping struct{ http.Handler }
type ephMapping struct{ http.Handler }

func newGameMapping(h http.Handler) mapping {
	g := &gameMapping{}
	g.Handler = http.StripPrefix(g.pathPrefix(), h)
	return g
}

func newStackMapping(h http.Handler) mapping {
	s := &stackMapping{}
	s.Handler = http.StripPrefix(s.pathPrefix(), h)
	return s
}

func newDrawingMapping(h http.Handler) mapping {
	d := &drawingMapping{}
	d.Handler = http.StripPrefix(d.pathPrefix(), h)
	return d
}

func newEphMapping(h http.Handler) mapping {
	e := &ephMapping{}
	e.Handler = http.StripPrefix(e.pathPrefix(), h)
	return e
}

func (g *gameMapping) pathPrefix() string {
	return "/game/"
}

func (g *gameMapping) canMap(obj interface{}) bool {
	switch obj.(type) {
	case model.Game:
		return true
	}
	return false
}

func (g *gameMapping) pathFor(obj interface{}) string {
	return g.pathPrefix() + obj.(model.Game).Gid() + "/"
}

func (s *stackMapping) pathPrefix() string {
	return "/stack/"
}

func (s *stackMapping) canMap(obj interface{}) bool {
	switch obj.(type) {
	case model.Stack:
		return true
	}
	return false
}

func (s *stackMapping) pathFor(obj interface{}) string {
	return s.pathPrefix() + obj.(model.Stack).Sid()
}

func (d *drawingMapping) pathPrefix() string {
	return "/drawing/"
}

func (d *drawingMapping) canMap(obj interface{}) bool {
	switch obj.(type) {
	case model.Drawing:
		return true
	}
	return false
}

func (d *drawingMapping) pathFor(obj interface{}) string {
	return d.pathPrefix() + obj.(model.Drawing).Did()
}

func (e *ephMapping) pathPrefix() string {
	return "/ephemeral/"
}

func (e *ephMapping) canMap(obj interface{}) bool {
	switch obj.(type) {
	case ephemeral.Ephemeris:
		return true
	}
	return false
}

func (e *ephMapping) pathFor(obj interface{}) string {
	return e.pathPrefix() + obj.(ephemeral.Ephemeris).Id()
}

type staticMapping struct {
	http.Handler
	prefix string
}

func newStaticMapping(dir, path string) *staticMapping {
	return &staticMapping{
		http.StripPrefix(path, http.FileServer(http.Dir(dir))),
		path}
}

func (s *staticMapping) pathPrefix() string {
	return s.prefix
}

func (s *staticMapping) canMap(interface{}) bool {
	return false
}

func (s *staticMapping) pathFor(interface{}) string {
	return ""
}
