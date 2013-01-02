package ephemeral

import (
	. "net/http"
	"net/url"
)

type Ephemera interface {
	EphemeralHandlerFor(id string) Handler
	Handler
}

type UserEphemera interface {
	Ephemera
	NewCreateUser(alias, email, pw string, dest *url.URL) interface{}
}

type Ephemeris interface {
	Id() string
	Handler
}

type EphemeraHandler struct {
	Ephemera
}
