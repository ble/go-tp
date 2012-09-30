package main

import (
  "log"
  "net/http"
  "fmt"
)

type AppRoutes struct {*Nexus}
func (_ AppRoutes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
  rr := AsRoute(*(r.URL))
  if rr != nil {
    u := rr.AsURL()
    fmt.Println((&u).String())
    rw.WriteHeader(200)
  } else {
    rw.WriteHeader(404)
  }
}

func main() {
  log.Fatal(http.ListenAndServe(":8080", AppRoutes{}))
}
