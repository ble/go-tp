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

func (e *ephemeraImpl) NewCreateUser(alias, email, pw string) interface{} {
	e.Lock()
	defer e.Unlock()
	id := e.NewId()
	base := ephemerisBase{id, e}
	eph := createUserImpl{alias, email, pw, base}
	e.ephemera[id] = eph
	return eph
}

type ephemerisBase struct {
	ephId string
	*ephemeraImpl
}

func (e ephemerisBase) Id() string {
	return e.ephId
}

type createUserImpl struct {
	alias, email, pw string
	ephemerisBase
}

func (n createUserImpl) ServeHTTP(w ResponseWriter, r *Request) {
	b := n.Backend
	newUser, err := b.CreateUser(n.email, n.alias, n.pw)
	if err == nil {
		n.removeEphemeralHandler(n.ephId)
		cookie := &Cookie{
			Name:     "userId",
			Value:    newUser.Uid(),
			Path:     "/",
			HttpOnly: true}
		SetCookie(w, cookie)
		w.WriteHeader(StatusOK)
	} else {
		Error(w, "", StatusInternalServerError)
	}
}
