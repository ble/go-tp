package switchboard

import (
	"ble/tpg/ephemeral"
	"ble/tpg/handler"
	"ble/tpg/persistence"
	"ble/tpg/room"
	"net/http"
	"net/url"
	"strings"
)

type Switchboard struct {
	mappings []mapping
	fallback http.Handler
	ephemera ephemeral.UserEphemera
}

func (s *Switchboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, m := range s.mappings {
		if strings.HasPrefix(path, m.pathPrefix()) {
			m.ServeHTTP(w, r)
			return
		}
	}
	s.fallback.ServeHTTP(w, r)
}

func (s *Switchboard) CanRoute(u url.URL) bool {
	path := u.Path
	for _, m := range s.mappings {
		if strings.HasPrefix(path, m.pathPrefix()) {
			return true
		}
	}
	return false
}

func (s *Switchboard) URLOf(i interface{}) *url.URL {
	for _, m := range s.mappings {
		if m.canMap(i) {
			loc, err := url.ParseRequestURI(m.pathFor(i))
			if err != nil {
				return nil
			}
			return loc
		}
	}
	return nil
}

func (s *Switchboard) GetEphemera() ephemeral.UserEphemera {
	return s.ephemera
}

func NewSwitchboard(b *persistence.Backend) *Switchboard {
	mappings := []mapping{nil, nil, nil, nil, nil}
	eph := ephemeral.NewEphemera(b)
	sb := &Switchboard{
		mappings: mappings,
		ephemera: eph,
		fallback: http.NotFoundHandler()}

	roomService := room.NewRoomService(sb, b)
	gameHandler := handler.NewGameHandler(roomService)
	mappings[0] = newGameMapping(gameHandler)

	stackHandler := handler.NewStackHandler(roomService)
	mappings[1] = newStackMapping(stackHandler)
	mappings[2] = newDrawingMapping(nil)
	mappings[3] = newEphMapping(eph)
	mappings[4] = newStaticMapping("./static/", "/static/")
	return sb
}
