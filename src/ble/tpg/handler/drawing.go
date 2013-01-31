package handler

import (
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

	playerId, _ := getPlayerId(r)
	userId, _ := getUserId(r)

	if isGet(r) {
		// the stack is complete ==>
		//        anyone can get the drawing.
		canGet := drawing.Stack().IsComplete()
		//        the drawing player can always get it
		canGet = canGet || drawing.Player().Pid() == playerId
		// the drawing is complete ==>
		//        any player who has completed a drawing
		//        in this stack can get it
		if !canGet && drawing.IsComplete() {
			stack := drawing.Stack()
			for _, otherDrawing := range stack.AllDrawings() {
				if otherDrawing.Player().Pid() == playerId &&
					otherDrawing.IsComplete() {
					canGet = true
					break
				}
			}
		}
		if canGet {
			jsonBytes, err := json.Marshal(drawing)
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
			drawing.Player().Pid() != playerId {
			Error(w, "not allowed to write to drawing", StatusBadRequest)
		} else {
			room.Draw(userId, playerId, drawing, r.Body)
			//TODO: process change to drawing
		}
	} else {
		Error(w, "method not allowed", StatusMethodNotAllowed)
	}
}
