package game

import (
	"encoding/json"
	"time"
)

func StartEvent() GameEvent {
	e := Event{EventType: "GameStart"}
	return GameEvent{time.Now(), e}
}

func FinishEvent() GameEvent {
	e := Event{EventType: "GameFinish"}
	return GameEvent{time.Now(), e}
}

func JoinEvent(name, id string) GameEvent {
	e := Event{EventType: "JoinGame", Name: name, Who: id}
	return GameEvent{time.Now(), e}
}

func PassEvent(fromId, toId, seqId string) GameEvent {
	e := Event{
		EventType: "PassSequence",
		FromWho:   fromId,
		ToWhom:    toId,
		What:      seqId,
	}
	return GameEvent{time.Now(), e}
}

type GameEvent struct {
	time.Time
	Payload interface{}
}

func (g GameEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Payload)
}

type Event struct {
	EventType string
	Name      string `json:",omitempty"`
	Who       string `json:",omitempty"`
	FromWho   string `json:",omitempty"`
	ToWhom    string `json:",omitempty"`
	What      string `json:",omitempty"`
	Error     string `json:",omitempty"`
}
