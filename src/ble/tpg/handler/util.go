package handler

import (
	. "net/http"
	"strings"
)

func isPost(r *Request) bool {
	return strings.ToLower(r.Method) == "post"
}

func isGet(r *Request) bool {
	return strings.ToLower(r.Method) == "get"
}

func pathParts(r *Request) []string {
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	if p[len(p)-1] == '/' {
		p = p[0 : len(p)-1]
	}
	return strings.Split(p, "/")
}

func getUserId(r *Request) (string, error) {
	uCookie, err := r.Cookie("userId")
	if err != nil {
		return "", err
	}
	return uCookie.Value, nil
}

func getPlayerId(r *Request) (string, error) {
	pCookie, err := r.Cookie("playerId")
	if err != nil {
		return "", err
	}
	return pCookie.Value, nil
}
