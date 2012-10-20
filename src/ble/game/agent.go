package game

import (
	. "ble/success"
	"errors"
	"time"
)

func NewGame() GameAgent {
	return newGameAgent(30*time.Second, 5*time.Second)
}

func newGameAgent(queueAge, tickPeriod time.Duration) GameAgent {
	chOut := make(chan GameEvent, 100)
	chIn := make(chan interface{})
	chQ := make(chan *eventQuery)

	agent := GameAgent{
		newGame(chOut),
		chOut,
		chIn,
		chQ,
		new(bool),
		new(bool),
		queueAge,
		tickPeriod}
	go agent.run()
	go agent.runEvents()
	return agent
}

type GameAgent struct {
	*Game
	events                     <-chan GameEvent
	messages                   chan interface{}
	queries                    chan *eventQuery
	gameStopped, eventsStopped *bool
	queueAge, tickPeriod       time.Duration
}

func (g GameAgent) GetGameEvents(last time.Time) ([]GameEvent, time.Time) {
	query := &eventQuery{last, make(chan []GameEvent)}
	g.queries <- query
	events := <-query.reply
	return events, query.lastQueried
}

func (g GameAgent) Shutdown() {
	close(g.Game.Events)
	close(g.messages)
}

func (g GameAgent) IsStarted() (bool, error) {
	msg := mIsStarted{make(Success), false}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.started, err
}
func (g GameAgent) Start() (bool, error) {
	msg := mStart{make(Success), false}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.started, err
}

func (g GameAgent) AddArtist(name string) (Artist, error) {
	msg := mAddArtist{make(Success), name, nil}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	if msg.created == nil {
		var artist0 Artist
		return artist0, errors.New("failed to add artist")
	}
	return *msg.created, err
}

func (g GameAgent) HasArtist(id string) (bool, error) {
	msg := mHasArtist{make(Success), id, false}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.present, err
}

func (g GameAgent) View() (interface{}, error) {
	msg := mView{make(Success), nil}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.view, err
}

func (g GameAgent) PassSequence(artistId string) (bool, error) {
	msg := mPassSequence{make(Success), artistId, false}
	g.messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.passed, err
}

//implementation: game state
func (g GameAgent) run() {
	for msg := range g.messages {
		switch m := msg.(type) {
		case *mIsStarted:
			m.started = g.Started
			m.Success <- true

		case *mStart:
			m.started = g.start()
			m.Success <- true

		case *mHasArtist:
			//ugly that this is so indirect
			m.present = g.NextArtist[m.id] != ""
			m.Success <- true

		case *mAddArtist:
			m.created = g.addArtist(m.name)
			m.Success <- true

		case *mView:
			m.view = g.viewJSON()
			m.Success <- true

		case *mPassSequence:
			m.passed = g.passSequence(m.artistId)
			m.Success <- true
		}
	}
	*g.gameStopped = true
}

//game state messages
type mIsStarted struct {
	Success
	started bool
}
type mStart struct {
	Success
	started bool
}
type mHasArtist struct {
	Success
	id      string
	present bool
}
type mAddArtist struct {
	Success
	name    string
	created *Artist
}
type mView struct {
	Success
	view interface{}
}
type mPassSequence struct {
	Success
	artistId string
	passed   bool
}

//implementation: event buffer
func (g GameAgent) runEvents() {
	var event GameEvent
	queue := make([]GameEvent, 0, 10)
	running := true
	ticks := time.Tick(g.tickPeriod)

	purgeEvents := func(events []GameEvent) []GameEvent {
		t := time.Now()
		eventsPrime := make([]GameEvent, 0, len(events))
		for _, v := range events {
			if t.Sub(v.Time) < g.queueAge {
				eventsPrime = append(eventsPrime, v)
			}
		}
		return eventsPrime
	}

	filterEvents := func(events []GameEvent, cutoff time.Time) []GameEvent {
		eventsPrime := make([]GameEvent, 0, len(events))
		for _, v := range events {
			if !v.Time.Before(cutoff) {
				eventsPrime = append(eventsPrime, v)
			}
		}
		return eventsPrime
	}

	for running {
		queue = purgeEvents(queue)
		select {
		case <-ticks:
			//loop again, purging events
		case event, running = <-g.events:
			queue = append(queue, event)
			if !running {
				*g.eventsStopped = true
			}
		case query := <-g.queries:
			query.reply <- filterEvents(queue, query.lastQueried)
			query.lastQueried = time.Now()
		}
	}
	*g.eventsStopped = true
}

type eventQuery struct {
	lastQueried time.Time
	reply       chan []GameEvent
}
