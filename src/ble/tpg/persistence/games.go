package persistence

import (
	"ble/hash"
	"ble/tpg/model"
	"database/sql"
	"errors"
	"sync"
)

type gamesBackend struct {
	*Backend
	createGame *sql.Stmt
	allGames   map[string]model.Game
	byName     map[string]model.Game
	*sync.RWMutex
}

func (g *gamesBackend) CreateGame(roomName string) (model.Game, error) {
	g.Lock()
	defer g.Unlock()
	if _, present := g.byName[roomName]; present {
		return nil, errors.New("game by that name already exists")
	}
	if err := g.prepStatement(
		"createGame",
		`INSERT INTO games
    (gid, started, complete, roomName)
    VALUES (?, 0, 0, ?);`,
		&g.createGame); err != nil {
		return nil, err
	}

	gameId := hash.EasyNonce(roomName)
	_, err := g.createGame.Exec(gameId, roomName)
	if err != nil {
		return nil, err
	}

	newGame := &game{
		g.gameBackend,
		roomName,
		gameId,
		make([]model.Player, 0, 0),
		make(map[model.Player]int),
		make(map[string]model.Player),
		make([]model.Stack, 0, 0),
		make(map[model.Player][]model.Stack),
		false,
		false}
	g.allGames[gameId] = newGame
	g.byName[roomName] = newGame
	return newGame, nil
}

func (g *gamesBackend) AllGames() map[string]model.Game {
	return g.allGames
}

func typecheckGames() model.Games {
	return &gamesBackend{}
}
