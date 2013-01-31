package persistence

import (
	"ble/hash"
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type gameBackend struct {
	*Backend
	markGameStarted, makeStack, markGameComplete, joinGame, passStack *sql.Stmt
	getStacks                                                         *sql.Stmt
}

type game struct {
	*gameBackend
	roomName     string
	gid          string
	players      []model.Player
	inverseOrder map[model.Player]int
	playersById  map[string]model.Player

	stacks       []model.Stack
	stacksInPlay map[model.Player][]model.Stack

	isComplete, isStarted bool
}

func (g *game) Gid() string {
	return g.gid
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
	playerId := hash.EasyNonce(u.Uid(), pseudonym, g.gid)
	if err := g.prepStatement(
		"joinGame",
		`INSERT INTO players
    (pid, pseudonym, gid, uid, playOrder)
    SELECT ? as pid,
           ? as pseudonym,
           ? as gid,
           ? as uid,
           count(ps.pid) as playOrder
    FROM players as ps
    WHERE ps.gid = ?;`,
		&g.joinGame); err != nil {
		return nil, err
	}
	_, err := g.joinGame.Exec(playerId, pseudonym, g.gid, u.Uid(), g.gid)
	if err != nil {
		return nil, err
	}

	player := &player{
		user:      u,
		pseudonym: pseudonym,
		pid:       playerId,
		game:      g,
	}
	g.inverseOrder[player] = len(g.players)
	g.players = append(g.players, player)
	g.playersById[player.Pid()] = player
	return player, nil
}

func (g *game) PlayerForId(pid string) model.Player {
	player, ok := g.playersById[pid]
	if !ok {
		return nil
	}
	return player
}

func (g *game) Stacks() []model.Stack {
	return g.stacks
}

func (g *game) StacksInProgress() map[model.Player][]model.Stack {
	return g.stacksInPlay
}

func (g *game) StacksFor(p model.Player) []model.Stack {
	return g.stacksInPlay[p]
}

func (g *game) PassStack(pFrom model.Player) (model.Stack, error) {
	//preconditions
	if g.IsComplete() {
		return nil, errors.New("can't pass a stack in a completed game!")
	}
	if !g.IsStarted() {
		return nil, errors.New("can't pass a stack in a game before it starts!")
	}
	if _, playerPresent := g.inverseOrder[pFrom]; !playerPresent {
		return nil, errors.New("that player isn't in this game!")
	}
	pFromStacks := g.stacksInPlay[pFrom]
	if len(pFromStacks) == 0 {
		return nil, errors.New("that player isn't holding any stacks!")
	}

	passedStack := pFromStacks[0]
	if !passedStack.TopDrawing().IsComplete() {
		return nil, errors.New("can't pass a stack with an incomplete drawing")
	}
	//figure out who'll be holding what after the pass...
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
		return nil, err
	}
	//execute statement
	_, err := g.passStack.Exec(pTo.Pid(), passedStack.Sid())
	if err != nil {
		return nil, err
	}
	//change in-memory structure
	g.stacksInPlay[pFrom] = pFromStacks

	//in-memory-only indication that no one is holding this stack any more...
	if !passedStack.IsComplete() {
		g.stacksInPlay[pTo] = pToStacks
	}
	return passedStack, nil
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
		"makeStack",
		`INSERT INTO stacks
    (sid, gid, complete, holdingPid)
    VALUES (?, ?, 0, ?);`,
		&g.makeStack); err != nil {
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

	if err := g.prepStatement(
		"addDrawing",
		`INSERT INTO drawings
    (did, sid, pid, stackOrder, complete)
    SELECT ? as did,
           ? as sid,
           ? as pid,
           ? as stackOrder,
           0 as complete
    FROM drawings as ds
    WHERE ds.sid = ?;`,
		&g.addDrawing); err != nil {
		return err
	}

	tx, err := g.Conn().Begin()
	if err != nil {
		return err
	}

	markGame := tx.Stmt(g.markGameStarted)
	makeStack := tx.Stmt(g.makeStack)
	addDrawing := tx.Stmt(g.addDrawing)

	if _, err := markGame.Exec(g.gid); err != nil {
		return err
		tx.Rollback()
	}
	N := len(g.players)
	stacks := make([]model.Stack, 0, N)
	stacksInPlay := make(map[model.Player][]model.Stack)
	for ix, player := range g.players {
		stackId := hash.EasyNonce(player.Pid(), ix)
		if _, err := makeStack.Exec(stackId, g.gid, player.Pid()); err != nil {
			tx.Rollback()
			return err
		}
		drawingId := hash.EasyNonce(player.Pid(), ix, stackId)
		if _, err := addDrawing.Exec(
			drawingId,
			stackId,
			player.Pid(),
			0,
			stackId); err != nil {
			tx.Rollback()
			return err
		}
		stack := &stack{
			stackBackend: g.stackBackend,
			sid:          stackId,
			g:            g,
			ds:           make([]model.Drawing, 0, N),
			completed:    false}
		g.recordNewStack(stack)
		drawing := &drawing{
			drawingBackend: g.drawingBackend,
			did:            drawingId,
			s:              stack,
			p:              player,
			content:        make([]interface{}, 0, 32),
			completed:      false}
		stack.ds = append(stack.ds, drawing)
		stacks = append(stacks, stack)
		stacksInPlay[player] = []model.Stack{stack}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	g.stacksInPlay = stacksInPlay
	g.stacks = stacks
	g.isStarted = true
	for _, stack := range g.stacks {
		for _, drawing := range stack.AllDrawings() {
			g.addDrawingToBackend(drawing)
		}
	}
	return nil
}
func typecheckGame() model.Game {
	return &game{}
}
