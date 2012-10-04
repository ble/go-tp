package main

import (
  "log"
  "net/http"
)

func main() {
  state := NewAppState()
  log.Fatal(http.ListenAndServe(":8080", state))
}
