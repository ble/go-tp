package ephemeral

import (
	. "net/http"
)

type Ephemera interface {
	EphemeralHandlerFor(id string) Handler
	Handler
}

type UserEphemera interface {
	Ephemera
	NewCreateUser(alias, email, pw string) interface{}
}

type Ephemeris interface {
	Id() string
	Handler
}

type EphemeraHandler struct {
	Ephemera
}
