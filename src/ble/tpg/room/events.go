package room

import (
	"encoding/json"
	"time"
)

type event struct {
	time.Time
	payload interface{}
}

type timeReq chan<- time.Time

type eventReq struct {
	events   chan<- []event
	lastTime time.Time
}

func (e event) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.payload)
}

const loopTime = 5 * time.Second

func (a *aRoom) processEvents() {
	ticks := time.NewTicker(loopTime)
	eventQueue := make([]event, 0, 100)
	select {
	case e := <-a.events:
		eventQueue = append(eventQueue, event{time.Now(), e})
	case r := <-a.eventRequests:
		if tReq, ok := r.(timeReq); ok {
			tReq <- time.Now()
		}
		if eReq, ok := r.(eventReq); ok {
			a.sendBackEvents(eReq, eventQueue)
		}
	case _ = <-ticks.C:
	}
	ticks.Stop()
}

func (a *aRoom) sendBackEvents(e eventReq, allEvents []event) {

}
