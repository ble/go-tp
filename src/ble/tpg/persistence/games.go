package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type gamesBackend struct {
	*Backend
	createGame *sql.Stmt
}

type games struct {
	*gamesBackend
	allGames map[string]model.Game
}

func (g *games) CreateGame(roomName string) (model.Game, error) {
	if _, present := g.allGames[roomName]; present {
		return nil, errors.New("game by that name already exists")
	}
	if err := g.prepStatement(
		"createGame",
		`INSERT INTO games
    (started, complete, roomName)
    VALUES (0, 0, ?);`,
		&g.createGame); err != nil {
		return nil, err
	}

	result, err := g.createGame.Exec(roomName)
	if err != nil {
		return nil, err
	}

	gid, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	newGame := &game{
		g.gameBackend,
		roomName,
		int(gid),
		make([]model.Player, 0, 0),
		make(map[model.Player]int),
		make(map[int]model.Player),
		make([]model.Stack, 0, 0),
		make(map[model.Player][]model.Stack),
		false,
		false}
	return newGame, nil
	return nil, nil
}

func (g *games) AllGames() map[string]model.Game {
	return g.allGames
}

func typecheckGames() model.Games {
	return &games{}
}
