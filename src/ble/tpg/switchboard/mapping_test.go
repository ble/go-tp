package switchboard

import (
	"ble/tpg/model"
	"ble/tpg/persistence"
	"os"
	"runtime/debug"
	"testing"
)

func dieOnErr(e error, t *testing.T) {
	if e != nil {
		debug.PrintStack()
		t.Fatal(e)
	}
}

func TestSwitchboardMapping(t *testing.T) {

	switchboard := NewSwitchboard()

	canRoundTrip := func(i interface{}) bool {
		u := switchboard.URLOf(i)
		if u == nil {
			return false
		}
		return switchboard.CanRoute(*u)
	}

	var backend *persistence.Backend
	var err error
	{
		backend, err = persistence.NewBackend("testdb")
		dieOnErr(err, t)
		defer os.Remove("testdb")

		backend.RegisterLogger(t)
		t.Log("set up backend")
	}

	var user1, user2 model.User
	{
		user1, err = backend.CreateUser("a@b.com", "fliffKnight", "paindeer")
		dieOnErr(err, t)
		user2, err = backend.CreateUser("q@b.com", "momoSandwich", "racelikeapisshorse")
		dieOnErr(err, t)
	}

	games := backend.CreateGamesService()
	var game model.Game
	var player1, player2 model.Player
	{
		game, err = games.CreateGame("grapnal vs. dognal")
		dieOnErr(err, t)
		player1, err = game.JoinGame(user1, "swaggerjacker")
		dieOnErr(err, t)
		player2, err = game.JoinGame(user2, "ailing vomit")
		dieOnErr(err, t)
	}
	player1.Pid()
	player2.Pid()
	err = game.Start()
	dieOnErr(err, t)
	t.Logf("Game URL: %s", switchboard.URLOf(game))
	if !canRoundTrip(game) {
		t.Fatal("can't roundtrip game")
	}

	stack1 := game.StacksFor(player1)[0]
	stack2 := game.StacksFor(player2)[0]
	t.Logf("Stack 1 URL: %s", switchboard.URLOf(stack1))
	t.Logf("Stack 2 URL: %s", switchboard.URLOf(stack2))
	if !canRoundTrip(stack1) {
		t.Fatal("can't roundtrip stack1")
	}
	if !canRoundTrip(stack2) {
		t.Fatal("can't roundtrip stack2")
	}
	drawing1 := stack1.TopDrawing()
	drawing2 := stack2.TopDrawing()
	t.Logf("Drawing 1 URL: %s", switchboard.URLOf(drawing1))
	t.Logf("Drawing 2 URL: %s", switchboard.URLOf(drawing2))
	if !canRoundTrip(drawing1) {
		t.Fatal("can't roundtrip drawing1")
	}
	if !canRoundTrip(drawing2) {
		t.Fatal("can't roundtrip drawing2")
	}

}
