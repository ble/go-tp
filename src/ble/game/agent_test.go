package game

import (
	"encoding/json"
	. "testing"
	. "time"
)

var time0 Time

func Test_Shutdown_Via_Close(t *T) {
	agent := NewGame()
	reachAround := agent.Game.Events
	reachAround <- GameEvent{Now(), "asdf"}

	agent.Start()

	_ = <-After(100 * Millisecond)

	started, err := agent.IsStarted()
	t.Log("Game is running: ", started)
	t.Log("started err", err)

	close(reachAround)
	close(agent.messages)
	_ = <-After(100 * Millisecond)

	if !*agent.gameStopped {
		t.Error("game has not stopped")
	}
	if !*agent.eventsStopped {
		t.Error("event queue has not stopped")
	}
}

func Test_Shutdown_Via_Interface(t *T) {
	agent := NewGame()
	agent.Shutdown()
	_ = <-After(100 * Millisecond)
	if !*agent.gameStopped {
		t.Error("game has not stopped")
	}
	if !*agent.eventsStopped {
		t.Error("event queue has not stopped")
	}

}

func Test_Getting_Game_Events(t *T) {
	agent := NewGame()

	a0, _ := agent.AddArtist("sammy")
	a1, _ := agent.AddArtist("buffalo vincent")

	agent.Start()

	_, afterErr := agent.AddArtist("shaniqua")
	if afterErr == nil {
		t.Error("artist joined after game started")
	} else {
		t.Log(afterErr)
	}

	agent.PassSequence(a0.Id)
	agent.PassSequence(a0.Id)
	agent.PassSequence(a1.Id)
	agent.PassSequence(a1.Id)
	agent.PassSequence(a0.Id)
	_ = <-After(100 * Millisecond)
	t.Log(time0)
	events, time1 := agent.GetGameEvents(time0)
	t.Log(time1)
	for _, v := range events {
		bytes, _ := json.Marshal(v)
		t.Log(string(bytes))
	}
	agent.Shutdown()
}

func Test_Event_Expiration(t *T) {
	agent := newGameAgent(325*Millisecond, 25*Millisecond)
	_, _ = agent.AddArtist("sammy")
	_ = <-After(80 * Millisecond)
	_, _ = agent.AddArtist("buffalo vincent")
	_ = <-After(80 * Millisecond)
	agent.Start()

	es1, _ := agent.GetGameEvents(time0)
	_ = <-After(200 * Millisecond)
	es2, _ := agent.GetGameEvents(time0)
	_ = <-After(100 * Millisecond)
	es3, _ := agent.GetGameEvents(time0)
	t.Log(es1)
	t.Log(es2)
	t.Log(es3)
	if len(es1) < len(es2) || len(es2) < len(es3) {
		t.Error("race-conditioned our way to nonsense land")
	}
	agent.Shutdown()
}

func Test_Event_Filtration(t *T) {

	agent := NewGame()
	agent.AddArtist("without spaces")
	_, time1 := agent.GetGameEvents(time0)
	es2, _ := agent.GetGameEvents(time1)
	agent.Start()
	es3, _ := agent.GetGameEvents(time1)
	if len(es2) != 0 {
		t.Error("failed to filter out event prior to query time")
	}
	if len(es3) != 1 {
		t.Error("failed to return single relevant event")
	}

	agent.Shutdown()
}
