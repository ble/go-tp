package handler

import (
	"ble/tpg/model"
	"ble/tpg/room"
	"encoding/json"
	. "net/http"
)

type drawingHandler struct {
	room.RoomService
}

func NewDrawingHandler(rs room.RoomService) Handler {
	return &drawingHandler{rs}
}

//GET <drawing-id>
//POST <drawing-id>
//POST <drawing-id>/complete

func (d *drawingHandler) ServeHTTP(w ResponseWriter, r *Request) {
	parts := pathParts(r)
	if len(parts) != 1 {
		NotFound(w, r)
		return
	}

	drawingId := parts[0]
	drawing, room, err := d.RoomService.GetDrawingAndRoom(drawingId)
	if err != nil {
		Error(w, err.Error(), StatusNotFound)
	}

	userId, _ := getUserId(r)
	sameUserId := func(uid string, d model.Drawing) bool {
		return uid == d.Player().User().Uid()
	}

	if isGet(r) {
		// the stack is complete ==>
		//        anyone can get the drawing.
		canGet := drawing.Stack().IsComplete()
		//        the drawing player can always get it
		canGet = canGet || sameUserId(userId, drawing)
		// the drawing is complete ==>
		//        any player who has completed a drawing
		//        in this stack can get it
		if !canGet && drawing.IsComplete() {
			stack := drawing.Stack()
			for _, otherDrawing := range stack.AllDrawings() {
				if sameUserId(userId, otherDrawing) &&
					otherDrawing.IsComplete() {
					canGet = true
					break
				}
			}
		}
		if canGet {
			jsonBytes, err := json.Marshal(drawing.Content())
			if err != nil {
				Error(w, err.Error(), StatusInternalServerError)
			} else {
				w.Write(jsonBytes)
			}
		} else {
			Error(w, "not allowed to read drawing", StatusForbidden)
		}
		return
	} else if isPost(r) {
		if drawing.IsComplete() ||
			!sameUserId(userId, drawing) {
			Error(w, "not allowed to write to drawing", StatusBadRequest)
		} else if len(parts) == 1 {
			err := room.Draw(userId, drawing, r.Body)
			if err != nil {
				Error(w, err.Error(), StatusBadRequest)
			}
		} else if len(parts) == 2 && parts[1] == "complete" {
			err := drawing.Complete()
			if err != nil {
				Error(w, err.Error(), StatusBadRequest)
			}
		} else {
			Error(w, "", StatusNotFound)
		}
	} else {
		Error(w, "method not allowed", StatusMethodNotAllowed)
	}
}
