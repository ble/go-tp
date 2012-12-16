package web

import (
	"net/http"
	"net/url"
)

type Switchboard interface {
	http.Handler

	// CanRoute returns true if the argument URL has the form of some URLs
	// understood by the Switchboard.  CanRoute returning false guarantees that
	// the URL does not correspond to a resource accessible through this
	// switchboard; CanRoute returning true does not guarantee that the URL does
	// correspond to a resource, but does mean that some portion of the URL is
	// "understandable" to the Switchboard.
	CanRoute(url.URL) bool

	// URLOf returns a URL at which a resource related (conflated?) with the
	// argument object may be accessed through this Switchboard.
	// The guarantee is that if URLOf(x) returns a non-nil URL, then there
	// exists some HTTP request with that URL such that if it is made to this
	// Switchboard, processing of that request will in some way access x.
	//
	// If the Switchboard does not allow access to the argument object, a
	// nil *URL will be returned.
	URLOf(interface{}) *url.URL
}
