package main

import (
	. "ble/drawing"
	"log"
	. "net/http"
)

func main() {
	Handle(
		"/static/",
		StripPrefix(
			"/static/",
			FileServer(Dir("./static"))))

	HandleFunc("/client", func(w ResponseWriter, r *Request) {
		ServeFile(w, r, "./static/drawing-client.html")
	})

	drawing := NewDrawingHandle()
	Handle("/", AsHandler(drawing))

	log.Fatal(ListenAndServe(":24769", nil))
}
