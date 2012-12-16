package persistence

import (
	DD "ble/drawing"
	"ble/tpg/model"
	"os"
	"runtime/debug"
	. "testing"
)

func helpErr(e error, t *T) {
	if e != nil {
		debug.PrintStack()
		t.Fatal(e)
	}
}

func TestCreateGame(t *T) {
	var backend *Backend
	var err error
	{
		backend, err = NewBackend("testdb")
		helpErr(err, t)
		defer os.Remove("testdb")

		backend.RegisterLogger(t)
		err = backend.createTables()
		helpErr(err, t)
		t.Log("set up backend")
	}

	var user1, user2, user3, userExtraneous model.User
	{
		userEmail := "the.bomb@thebomb.com"
		userAlias := "scatman juan"
		userPw := "asdfff"
		user1, err = backend.CreateUser(userEmail, userAlias, userPw)
		helpErr(err, t)
		user1, err = backend.LogInUser(userAlias, userPw)
		helpErr(err, t)
		if user1 == nil {
			t.Fatal("didn't log user in")
		} else {
			t.Log("created user and logged in")
		}
		user2, err = backend.CreateUser("wade@boggs", "biff", "qrstu")
		helpErr(err, t)

		user3, err = backend.CreateUser("b@q", "snafflebirt", "qrstu")
		helpErr(err, t)

		userExtraneous, err = backend.CreateUser("oafi", "gambino", "qrtsu")
		helpErr(err, t)

		duplicateEmail, err := backend.CreateUser(
			"the.bomb@thebomb.com",
			"swaggerjacker",
			"00asdf")
		if duplicateEmail != nil || err == nil {
			t.Fatal("allowed duplicate email")
		}

		duplicateAlias, err := backend.CreateUser(
			"fongle",
			user3.Alias(),
			"qrtsu")
		if duplicateAlias != nil || err == nil {
			t.Fatal("allowed duplicate alias")
		}
		t.Log("created players")
	}

	var gamesService model.Games
	var game model.Game
	{
		gamesService = &games{backend.gamesBackend, make(map[string]model.Game)}
		gameName := "grapnal vs. dognel"
		game, err = gamesService.CreateGame(gameName)
		helpErr(err, t)
		noGame, err := gamesService.CreateGame(gameName)
		if noGame != nil || err == nil {
			t.Fatal("allowed a game by a duplicate name")
		}
		t.Log("created a game")
	}

	var player1, player2, player3 model.Player
	{
		p1Name, p2Name, p3Name := "wizrad", "P-knee Sir Prize", "grapnal"
		player1, err = game.JoinGame(user1, p1Name)
		helpErr(err, t)
		player2, err = game.JoinGame(user2, p2Name)
		helpErr(err, t)
		player3, err = game.JoinGame(user3, p3Name)
		helpErr(err, t)
		noPlayer, err := game.JoinGame(userExtraneous, p1Name)
		if noPlayer != nil || err == nil {
			t.Fatal("allowed a duplicate player name")
		}
		noPlayer, err = game.JoinGame(user1, "not a sockpuppet")
		if noPlayer != nil || err == nil {
			t.Fatal("allowed one user to join the game twice")
		}
	}

	{
		err = game.Complete()
		if err == nil {
			t.Fatal("completed game before it started")
		}
		err = game.Start()
		helpErr(err, t)
		err = game.Start()
		if err == nil {
			t.Fatal("started game twice")
		}
		t.Log("started the game")
	}
	{
		if len(game.Stacks()) != 3 {
			t.Fatal("didn't create one stack per player")
		}

		stacks0 := game.StacksInProgress()[player1]
		if len(stacks0) != 1 {
			t.Fatal("player not holding a single stack")
		}
		theStack := stacks0[0]
		if len(theStack.AllDrawings()) != 1 {
			t.Fatal("the stack does not have a single drawing")
		}
		theDrawing := theStack.TopDrawing()
		helpErr(theDrawing.Add(DD.DefaultDrawPart), t)
		helpErr(theDrawing.Add(DD.DefaultDrawPart), t)
		helpErr(theDrawing.Complete(), t)
		helpErr(game.PassStack(player1), t)
		helpErr(game.PassStack(player2), t)
		player3.Pid()
	}
	/*
		_, err = theStack.AddDrawing(player2)
		helpErr(err, t)
		helpErr(game.PassStack(player3), t)
	*/
	/*
		drawing3, err := theStack.AddDrawing(p)
		helpErr(err, t)
		t.Log(err)
		t.Log(len(theStack.AllDrawings()))
		t.Logf("%#v", drawing2)
		t.Logf("%#v", drawing3)
	*/
}
