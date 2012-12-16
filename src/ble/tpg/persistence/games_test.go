package persistence

import (
	DD "ble/drawing"
	"ble/tpg/model"
	"os"
	"runtime/debug"
	. "testing"
)

func dieOnErr(e error, t *T) {
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
		dieOnErr(err, t)
		defer os.Remove("testdb")

		backend.RegisterLogger(t)
		err = backend.createTables()
		dieOnErr(err, t)
		t.Log("set up backend")
	}

	var user1, user2, user3, userExtraneous model.User
	{
		userEmail := "the.bomb@thebomb.com"
		userAlias := "scatman juan"
		userPw := "asdfff"
		user1, err = backend.CreateUser(userEmail, userAlias, userPw)
		dieOnErr(err, t)
		user1, err = backend.LogInUser(userAlias, userPw)
		dieOnErr(err, t)
		if user1 == nil {
			t.Fatal("didn't log user in")
		} else {
			t.Log("created user and logged in")
		}
		user2, err = backend.CreateUser("wade@boggs", "biff", "qrstu")
		dieOnErr(err, t)

		user3, err = backend.CreateUser("b@q", "snafflebirt", "qrstu")
		dieOnErr(err, t)

		userExtraneous, err = backend.CreateUser("oafi", "gambino", "qrtsu")
		dieOnErr(err, t)

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
		dieOnErr(err, t)
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
		dieOnErr(err, t)
		player2, err = game.JoinGame(user2, p2Name)
		dieOnErr(err, t)
		player3, err = game.JoinGame(user3, p3Name)
		dieOnErr(err, t)
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
		dieOnErr(err, t)
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
		dieOnErr(theDrawing.Add(DD.DefaultDrawPart), t)
		dieOnErr(theDrawing.Add(DD.DefaultDrawPart), t)
		dieOnErr(theDrawing.Complete(), t)
		dieOnErr(game.PassStack(player1), t)
		_, err := theStack.AddDrawing(game.NextPlayer(player1))
		dieOnErr(err, t)
		errPassIncomplete := game.PassStack(player2)
		if errPassIncomplete == nil {
			t.Fatal("player passed stack w/ incomplete drawing")
		}

		if len(game.StacksInProgress()[player1]) != 0 {
			t.Fatal("player 1 still holding stacks despite passing all of them")
		}
		t.Log("player 1 passed 1 stack")

		for i := 0; i < 2; i++ {
			stacks1 := game.StacksInProgress()[player2]
			theStack = stacks1[0]
			dieOnErr(theStack.TopDrawing().Complete(), t)
			dieOnErr(game.PassStack(player2), t)
			_, err := theStack.AddDrawing(game.NextPlayer(player2))
			dieOnErr(err, t)
		}

		if len(game.StacksInProgress()[player2]) != 0 {
			t.Fatal("player 2 still holding stacks despite passing all of them")
		}
		t.Log("player 2 passed 2 stacks")

		stacks2 := game.StacksInProgress()[player3]
		if len(stacks2) != 3 {
			t.Fatal("player 3 is not holding all of the stacks")
		}
		t.Log("player 3 holding 3 stacks")
	}
}
