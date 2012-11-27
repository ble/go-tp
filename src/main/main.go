package main

import (
	. "ble/drawing"
	"log"
	. "net/http"
)

func main() {
	drawing := NewDrawingHandle()
	Handle("/", AsHandler(drawing))
	HandleFunc("/client", func(w ResponseWriter, r *Request) {
		ServeFile(w, r, "./static/drawing-client.html")
	})
	//TODO: implement static resources other than client
	//with http.StripPrefix, http.Dir
	HandleFunc("/drawing-ui.js", func(w ResponseWriter, r *Request) {
		ServeFile(w, r, "./static/drawing-ui.js")
	})
	log.Fatal(ListenAndServe(":24769", nil))
}
