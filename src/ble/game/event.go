package game

import (
	"encoding/json"
	"time"
)

func startEvent() GameEvent {
	e := event{EventType: "GameStart"}
	return GameEvent{time.Now(), e}
}

func finishEvent() GameEvent {
	e := event{EventType: "GameFinish"}
	return GameEvent{time.Now(), e}
}

func joinEvent(name, id string) GameEvent {
	e := event{EventType: "JoinGame", Name: name, Who: id}
	return GameEvent{time.Now(), e}
}

func passEvent(fromId, toId, seqId string) GameEvent {
	e := event{
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

type event struct {
	EventType string
	Name      string `json:",omitempty"`
	Who       string `json:",omitempty"`
	FromWho   string `json:",omitempty"`
	ToWhom    string `json:",omitempty"`
	What      string `json:",omitempty"`
}
