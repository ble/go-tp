package room

import (
	"encoding/json"
	"time"
)

func (a *aRoom) GetEvents(uid, pid string, lastQuery Time) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

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
const filterAge = 120 * time.Second

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
		tNow := time.Now()
		var ix int
		var ev event
		for ix, ev = range eventQueue {
			if tNow.Sub(ev.Time) < filterAge ||
				ix == len(eventQueue)-1 {
				break
			}
		}
		eventQueue = eventQueue[ix:]
	}
	ticks.Stop()
}

func (a *aRoom) sendBackEvents(e eventReq, allEvents []event) {
	cutoff := e.lastTime
	count := 0
	for _, ev := range allEvents {
		if ev.Time.After(cutoff) {
			count++
		}
	}
	toSendBack := make([]event, count, count)
	count = 0
	for _, ev := range allEvents {
		if ev.Time.After(cutoff) {
			toSendBack[count] = ev
			count++
		}
	}
	e.events <- toSendBack
}
