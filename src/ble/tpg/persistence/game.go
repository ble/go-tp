package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type gameBackend struct {
	*Backend
	markGameStarted, makeStacks, markGameComplete, joinGame, passStack *sql.Stmt
	getStacks                                                          *sql.Stmt
}

type game struct {
	*gameBackend
	roomName     string
	gid          int
	players      []model.Player
	inverseOrder map[model.Player]int
	playersById  map[int]model.Player

	stacks       []model.Stack
	stacksInPlay map[model.Player][]model.Stack

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
		return nil, errors.New("Sorry, you can't join a game that's already started.")
	}
	if g.IsComplete() {
		return nil, errors.New("Sorry, you can't join a completed game.")
	}
	if err := g.prepStatement(
		"joinGame",
		`INSERT INTO players
    (pseudonym, gid, uid, playOrder)
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
	g.playersById[player.Pid()] = player
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
	if g.IsComplete() {
		return errors.New("can't pass a stack in a completed game!")
	}
	if !g.IsStarted() {
		return errors.New("can't pass a stack in a game before it starts!")
	}
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
	if err != nil {
		return err
	}
	//change in-memory structure
	g.stacksInPlay[pFrom] = pFromStacks

	//in-memory-only indication that no one is holding this stack any more...
	if !passedStack.IsComplete() {
		g.stacksInPlay[pTo] = pToStacks
	}
	return nil
}

func (g *game) IsComplete() bool {
	return g.isComplete
}

func (g *game) IsStarted() bool {
	return g.isStarted
}

func (g *game) Complete() error {
	if g.IsComplete() {
		return errors.New("game already completed!")
	}
	if !g.IsStarted() {
		return errors.New("game not yet started!")
	}
	//check to see if all stacks are complete
	for _, stack := range g.stacks {
		if !stack.IsComplete() {
			return errors.New("Not all stacks are complete!")
		}
	}

	if err := g.prepStatement(
		"markGameComplete",
		`UPDATE games
  SET complete = 1
  WHERE gid = ?;`,
		&g.markGameComplete); err != nil {
		return err
	}

	_, err := g.markGameComplete.Exec(g.gid)
	if err != nil {
		return err
	}

	g.isComplete = true
	return nil
}

func (g *game) Start() error {
	if g.IsStarted() {
		return errors.New("game already started!")
	}
	if g.IsComplete() {
		return errors.New("can't start a complete game!")
	}

	if err := g.prepStatement(
		"markGameStarted",
		`UPDATE games
    SET started = 1
    WHERE gid = ?`,
		&g.markGameStarted); err != nil {
		return err
	}

	if err := g.prepStatement(
		"makeStacks",
		`INSERT INTO stacks
    (gid, complete, holdingPid)
    SELECT ? as gid,
           0 as complete,
           ps.pid as holdingPid
    FROM players as ps
    WHERE ps.gid = ?;`,
		&g.makeStacks); err != nil {
		return err
	}

	if err := g.prepStatement(
		"getStacks",
		`SELECT sid, holdingPid
    FROM stacks
    WHERE gid = ?`,
		&g.getStacks); err != nil {
		return err
	}

	tx, err := g.Conn().Begin()
	if err != nil {
		return err
	}

	markGame := tx.Stmt(g.markGameStarted)
	makeStacks := tx.Stmt(g.makeStacks)
	getStacks := tx.Stmt(g.getStacks)

	if _, err := markGame.Exec(g.gid); err != nil {
		return err
		tx.Rollback()
	}

	if _, err := makeStacks.Exec(g.gid, g.gid); err != nil {
		return err
		tx.Rollback()
	}

	rows, err := getStacks.Query(g.gid)
	if err != nil {
		return err
		tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	g.stacks = make([]model.Stack, len(g.players), len(g.players))

	for i := 0; rows.Next(); i++ {
		var sid, holdingPid int
		err := rows.Scan(&sid, &holdingPid)
		if err != nil {
			return err
		}
		startingPlayer := g.playersById[holdingPid]
		theStack := &stack{
			g.stackBackend,
			sid,
			g,
			make([]model.Drawing, 0, 0),
			false}
		g.stacks[i] = theStack
		g.stacksInPlay[startingPlayer] = []model.Stack{theStack}
		if _, err := theStack.AddDrawing(startingPlayer); err != nil {
			return err
		}
	}
	g.isStarted = true
	return nil
}
func typecheckGame() model.Game {
	return &game{}
}
