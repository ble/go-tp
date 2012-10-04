package main

import (
  "net/url"
  "net/http"
)


type Route interface {
  Matches(url url.URL) bool
  AsURL() url.URL
  CanHandle(request *http.Request) bool
  ServeHTTP(http.ResponseWriter, *http.Request)
}

func GetRoute(r *http.Request, a *AppState) Route {
  routes := []Route{&RouteRoot{state: a}, &RouteCreate{state: a}}
  for _, route := range routes {
    if route.CanHandle(r) {
      return route
    }
  }
  return nil    
}


