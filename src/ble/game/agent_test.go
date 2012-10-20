package game

import (
	"encoding/json"
	. "testing"
	. "time"
)

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
	var time0 Time
	t.Log(time0)
	events, time1 := agent.GetGameEvents(time0)
	t.Log(time1)
	for _, v := range events {
		bytes, _ := json.Marshal(v)
		t.Log(string(bytes))
	}
}
