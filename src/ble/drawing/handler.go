package drawing

import (
	"encoding/json"
	. "net/http"
	"strings"
)

type drawingHandler struct {
	DrawingHandle
}

func (d drawingHandler) ServeHTTP(w ResponseWriter, r *Request) {
	method := strings.ToUpper(r.Method)
	if method == "GET" {
		d.ServeGet(w, r)
	} else if method == "POST" {
		d.ServePost(w, r)
	} else {
		w.WriteHeader(StatusMethodNotAllowed)
	}
}

func (d drawingHandler) ServeGet(w ResponseWriter, r *Request) {
	bytes, err := d.Read()
	if err != nil {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte("Internal server error: "))
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(StatusOK)
		w.Write(bytes)
	}
}

func (d drawingHandler) ServePost(w ResponseWriter, r *Request) {
	//reject overly-long post
	if r.ContentLength >= 1024*1024*16 {
		w.WriteHeader(StatusRequestEntityTooLarge)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "" {
		w.WriteHeader(StatusUnsupportedMediaType)
		w.Write([]byte("only json accepted"))
		return
	}

	posted := new(DrawPart)
	err := json.NewDecoder(r.Body).Decode(posted)
	if err != nil {
		w.WriteHeader(StatusBadRequest)
		w.Write([]byte("json failed to parse"))
		return
	}

	err = d.Draw(*posted)
	if err != nil {
		w.WriteHeader(StatusInternalServerError)
		w.Write([]byte("Internal server error: "))
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(StatusOK)
		json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
		}{"ok"})
		return
	}
}
