package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type RouteCreate struct {
	state AppHandle
	info  *UserInfo
}

func (_ *RouteCreate) Matches(url url.URL) bool {
	return url.Path == "/create" || url.Path == "create"
}

func (_ *RouteCreate) AsURL() url.URL {
	return url.URL{Path: "/create"}
}

func (route *RouteCreate) CanHandle(request *http.Request) bool {
	if route.Matches(*request.URL) {
		route.info = new(UserInfo)
		route.info.ExtractFrom(request)

		return true
	}
	return false
}

func (route *RouteCreate) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if strings.ToLower(req.Method) != "post" {
		http.Error(resp, "", 405)
		return
	}

	if req.Header["Content-Type"][0] != "application/json" {
		http.Error(resp, "", 406)
		return
	}

	bodyJson := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&bodyJson)
	if error != nil {
		http.Error(resp, "", 500)
		return
	}

	info := route.info
	state := route.state
	if info == nil {
		http.Error(resp, "", 500)
		return
	}
	if can, _ := state.UserCanCreateRoom(info.UserId); can {
		resp.WriteHeader(200)
		/*roomId :=*/ state.CreateRoom(req, info.UserId, bodyJson)
		return
	} else {

		resp.WriteHeader(401)
		return
	}

}
