package switchboard

import (
	"net/http"
)

type mapping interface {
	pathPrefix() string
	canMap(interface{}) bool
	pathFor(interface{}) string
	http.Handler
}
