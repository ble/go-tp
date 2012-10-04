package main

import (
"net/url"
"net/http"
"strings"
)

type RouteRoot struct {
  state *AppState
  info *UserInfo
}

func (_ *RouteRoot) Matches(url url.URL) bool {
  return url.Path == "/" || url.Path == ""
}

func (_ *RouteRoot) AsURL() url.URL {
  return url.URL{Path: "/"}
}

func (route *RouteRoot) CanHandle(request *http.Request) bool {
  if route.Matches(*request.URL) {
    route.info = new(UserInfo)
    route.info.ExtractFrom(request)
    
    return true
  }
  return false
}

func (route *RouteRoot) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
  if strings.ToLower(req.Method) != "get" {
    http.Error(resp, "", 405)
    return
  }

  info := route.info
  state := route.state
  if info == nil {
    http.Error(resp, "", 500)
    return
  }
  
  //no user id assigned
  if info.UserId == "" {
    userId := state.CreateUserId(req)
    cookie := http.Cookie{
                Name: "user",
                Value: userId,
                Path: route.AsURL().Path,
                HttpOnly: true}
    http.SetCookie(resp, &cookie)
    resp.WriteHeader(200)
    return
  }

  if state.UserExists(info.UserId) {
    resp.WriteHeader(200)
    return
  } else {
    route.info.UserId = ""
    route.ServeHTTP(resp, req)
  }

  
}


