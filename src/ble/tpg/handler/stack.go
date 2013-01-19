package handler

import (
	"ble/tpg/room"
	"encoding/json"
	. "net/http"
)

type stackHandler struct {
	room.RoomService
}

func NewStackHandler(rs room.RoomService) Handler {
	return &stackHandler{rs}
}

func (s *stackHandler) ServeHTTP(w ResponseWriter, r *Request) {
	parts := pathParts(r)
	if len(parts) != 1 {
		NotFound(w, r)
		return
	}

	stackId := parts[0]
	stack, _, err := s.RoomService.GetStackAndRoom(stackId)
	if err != nil {
		Error(w, err.Error(), StatusNotFound)
		return
	}
	obj := room.StackJson{stack, s.RoomService.GetSwitchboard()}
	stackJson, _ := json.Marshal(obj)
	w.Write(stackJson)

}
