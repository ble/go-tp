package handler

import (
	. "net/http"
)

type handlerJoinClient struct{}

func (handlerJoinClient) ServeHTTP(w ResponseWriter, r *Request) {
	ServeFile(w, r, "./static/html/join-client.html")
}
