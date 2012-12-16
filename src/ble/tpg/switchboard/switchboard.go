package switchboard

import (
	"ble/web"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type switchboard struct {
	mappings []mapping
	fallback http.Handler
}

func (s *switchboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, m := range s.mappings {
		if strings.HasPrefix(path, m.pathPrefix()) {
			m.ServeHTTP(w, r)
			return
		}
	}
	s.fallback.ServeHTTP(w, r)
}

func (s *switchboard) CanRoute(u url.URL) bool {
	path := u.Path
	for _, m := range s.mappings {
		if strings.HasPrefix(path, m.pathPrefix()) {
			return true
		}
	}
	return false
}

func (s *switchboard) URLOf(i interface{}) (*url.URL, error) {
	for _, m := range s.mappings {
		if m.canMap(i) {
			loc, err := url.ParseRequestURI(m.pathFor(i))
			if err != nil {
				return nil, err
			}
			return loc, nil
		}
	}
	return nil, fmt.Errorf("don't know how to map %#v", i)
}

func NewSwitchboard() web.Switchboard {
	return &switchboard{
		mappings: []mapping{
			newGameMapping(nil),
			newStackMapping(nil),
			newDrawingMapping(nil)}}
}
