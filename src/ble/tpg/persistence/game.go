package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type gameBackend struct {
	*Backend
	joinGame, passStack *sql.Stmt
}

type game struct {
	*gameBackend
	gid          int
	players      []model.Player
	inverseOrder map[model.Player]int

	stacks         []model.Stack
	stacksInPlay   map[model.Player][]model.Stack
	stacksFinished map[model.Stack]bool

	isComplete, isStarted bool
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
	if g.IsStarted() {
		return errors.New("Sorry, you can't join a game that's already started.")
	}
	if g.IsComplete() {
		return errors.New("Sorry, you can't join a completed game.")
	}
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
	return g.stacks
}

func (g *game) StacksInProgress() map[model.Player][]model.Stack {
	return g.stacksInPlay
}

func (g *game) PassStack(pFrom model.Player) error {
	//preconditions
	//  if !//started and completed condition here
	if _, playerPresent := g.inverseOrder[pFrom]; !playerPresent {
		return errors.New("that player isn't in this game!")
	}
	pFromStacks := g.stacksInPlay[pFrom]
	if len(pFromStacks) == 0 {
		return errors.New("that player isn't holding any stacks!")
	}

	//figure out who'll be holding what after the pass...
	passedStack := pFromStacks[0]
	pFromStacks = pFromStacks[1:]
	pTo := g.NextPlayer(pFrom)
	pToStacks := g.stacksInPlay[pTo]
	pToStacks = append(pToStacks, passedStack)

	//prep statement
	if err := g.prepStatement(
		"passStack",
		`UPDATE stacks
  SET holdingPid = ?
  WHERE sid = ?;`,
		&g.passStack); err != nil {
		return err
	}
	//execute statement
	_, err := g.passStack.Exec(pTo.Pid(), passedStack.Sid())
	//change in-memory structure
	g.stacksInPlay[pFrom] = pFromStacks
	g.stacksInPlay[pTo] = pToStacks
}

func (g *game) IsComplete() bool {
	return g.isComplete
}

func (g *game) IsStarted() bool {
	return g.isStarted
}

func (g *game) Complete() error {
	//if already complete or not started, that's an error
	//tx time!
	//mark this as complete
	//mark all stacks as complete
	//remove all stacks from play
	//update in-memory strucutres
}

func (g *game) Start() error {
	//if already started or completed, that's an error
	//tx time
	//mark this as started
	//create a stack for each player
	//create a drawing in each stack
}
func typecheckGame() model.Game {
	return &game{}
}
