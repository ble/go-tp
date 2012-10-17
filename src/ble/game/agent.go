package game

import . "ble/success"

func NewGame() GameAgent {
	chOut := make(chan GameEvent)
	chIn := make(chan interface{})

	agent := GameAgent{newGame(chOut), chOut, chIn}
	go agent.run(chIn)
	return agent
}

type GameAgent struct {
	*Game
	GameEvents <-chan GameEvent
	Messages   chan<- interface{}
}

func (g GameAgent) IsStarted() (bool, error) {
	msg := mStart{make(Success), false}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.started, err
}
func (g GameAgent) Start() (bool, error) {
	msg := mStart{make(Success), false}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.started, err
}

func (g GameAgent) AddArtist(name string) (Artist, error) {
	msg := mAddArtist{make(Success), name, nil}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return *msg.created, err
}

func (g GameAgent) HasArtist(id string) (bool, error) {
	msg := mHasArtist{make(Success), id, false}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.present, err
}
func (g GameAgent) View() (interface{}, error) {
	msg := mView{make(Success), nil}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.view, err
}

func (g GameAgent) PassSequence(artistId string) (bool, error) {
	msg := mPassSequence{make(Success), artistId, false}
	g.Messages <- &msg
	err := msg.SucceededIn(Second)
	return msg.passed, err
}

//implementation
func (g GameAgent) run(messages <-chan interface{}) {
	for msg := range messages {
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
}

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
