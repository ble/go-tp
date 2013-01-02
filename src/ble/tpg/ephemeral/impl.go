package ephemeral

import (
	"ble/tpg/persistence"
	. "net/http"
	"path"
	"sync"
)

func NewEphemera(b *persistence.Backend) UserEphemera {
	return &ephemeraImpl{
		sync.Mutex{},
		b,
		make(map[string]Ephemeris),
		&trivialIdGenerator{}}
}

type ephemeraImpl struct {
	sync.Mutex
	*persistence.Backend
	ephemera map[string]Ephemeris
	IdGenerator
}

func (e *ephemeraImpl) ServeHTTP(w ResponseWriter, r *Request) {
	reqPath := r.URL.Path
	for path.Dir(reqPath) != "." {
		reqPath = path.Dir(reqPath)
	}
	e.EphemeralHandlerFor(reqPath).ServeHTTP(w, r)
}

func (e *ephemeraImpl) EphemeralHandlerFor(id string) Handler {
	e.Lock()
	defer e.Unlock()
	if ephemeris, present := e.ephemera[id]; present {
		return ephemeris
	}
	return NotFoundHandler()
}

func (e *ephemeraImpl) removeEphemeralHandler(id string) {
	e.Lock()
	defer e.Unlock()
	delete(e.ephemera, id)
}

type ephemerisBase struct {
	ephId string
	*ephemeraImpl
}

func (e ephemerisBase) Id() string {
	return e.ephId
}
