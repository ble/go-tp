package ephemeral

import (
	. "net/http"
	"net/url"
)

func (e *ephemeraImpl) NewCreateUser(
	alias,
	email,
	pw string,
	destination *url.URL) interface{} {
	e.Lock()
	defer e.Unlock()
	id := e.NewId()
	base := ephemerisBase{id, e}
	eph := createUserImpl{alias, email, pw, destination, base}
	e.ephemera[id] = eph
	return eph
}

type createUserImpl struct {
	alias, email, pw string
	dest             *url.URL
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
		w.Header().Add("Location", n.dest.String())
		w.WriteHeader(StatusSeeOther)
	} else {
		Error(w, "", StatusInternalServerError)
	}
}
