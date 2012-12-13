package persistence

import (
	"ble/tpg/model"
	"database/sql"
)

type gameBackend struct {
	*Backend
	joinGame *sql.Stmt
}

type game struct {
	*gameBackend
	gid          int
	players      []model.Player
	inverseOrder map[model.Player]int

	stacks         []model.Stack
	stacksInPlay   map[model.Player][]model.Stack
	stacksFinished map[model.Stack]bool
}

func (g *game) Players() []model.Player {
	return g.players
}

func (g *game) NextPlayer(p model.Player) model.Player {
	index, present := g.inverseOrder[p]
	if !present {
		return nil
	}
	index = (index + 1) % len(g.players)
	return g.players[index]
}

func (g *game) JoinGame(u model.User, pseudonym string) (model.Player, error) {
	if err := g.prepStatement(
		"joinGame",
		`INSERT INTO players
    SELECT ? as pseudonym,
           ? as gid,
           ? as uid,
           count(ps.pid) as playOrder
    FROM players as ps
    WHERE ps.gid = ?;`,
		&g.joinGame); err != nil {
		return nil, err
	}
	result, err := g.joinGame.Exec(pseudonym, g.gid, u.Uid(), g.gid)
	if err != nil {
		return nil, err
	}
	playerId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	player := &player{
		user:      u,
		pseudonym: pseudonym,
		pid:       int(playerId),
		game:      g,
	}
	g.inverseOrder[player] = len(g.players)
	g.players = append(g.players, player)
	return player, nil
}

func (g *game) Stacks() []model.Stack {

}

func typecheckGame() model.Game {
	return &game{}
}
