package room

import (
	"encoding/json"
	"errors"
	"time"
)

func (a *aRoom) GetEvents(uid, pid string, lastQuery time.Time) (interface{}, error) {
	if a.game.PlayerForId(pid) == nil {
		return nil, errors.New("no such player")
	}
	reqChan := make(chan eventResponse)
	req := eventReq{reqChan, lastQuery}
	a.eventRequests <- req
	resp := <-reqChan
	return resp, nil
}

func (a *aRoom) GetLastEventTime() time.Time {
	reqChan := make(chan eventResponse)
	req := eventReq{reqChan, time.Now()}
	a.timeReqs <- req
	resp := <-reqChan
	return resp.LastTime
}

type event struct {
	time.Time
	payload interface{}
}

func (e event) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.payload)
}

type timeReq chan<- time.Time

type eventReq struct {
	response chan<- eventResponse
	lastTime time.Time
}

type eventResponse struct {
	LastTime time.Time
	Events   []event
}

func (eR eventResponse) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["events"] = eR.Events
	result["lastTime"] = eR.LastTime.UnixNano() / 1000
	return json.Marshal(result)
}

const loopTime = 5 * time.Second
const filterAge = 120 * time.Second

func (a *aRoom) processEvents() {
	ticks := time.NewTicker(loopTime)
	eventQueue := make([]event, 0, 100)
	for {
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
		case t := <-a.timeReqs:
			t.response <- eventResponse{time.Now(), make([]event, 0, 0)}
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
	}
	ticks.Stop()
}

func (a *aRoom) sendBackEvents(e eventReq, allEvents []event) {
	cutoff := e.lastTime
	now := time.Now()
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
	e.response <- eventResponse{now, toSendBack}
}
