package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type stackBackend struct {
	*Backend
	addDrawing, completeStack, wipeOrder *sql.Stmt
}

func (sb *stackBackend) validateStack(s model.Stack) bool {
	return true
}

type stack struct {
	*stackBackend
	sid       int
	g         model.Game
	ds        []model.Drawing
	completed bool
}

func typeCheckStack() model.Stack {
	return &stack{}
}

func (s *stack) Game() model.Game {
	return s.g
}

func (s *stack) AllDrawings() []model.Drawing {
	return s.ds
}

func (s *stack) AddDrawing(p model.Player) (model.Drawing, error) {
	if err := s.prepStatement(
		"addDrawing",
		`INSERT INTO drawings (sid, pid, stackOrder, complete)
    SELECT ? as sid,
           ? as pid,
           count(ds.pid) as stackOrder,
           false as complete
    FROM drawings as ds;`,
		&s.addDrawing); err != nil {
		return nil, err
	}
	result, err := s.addDrawing.Exec(s.sid, p.Pid())
	if err != nil {
		return nil, err
	}
	drawingId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	drawing := &drawing{
		drawingBackend: s.drawingBackend,
		did:            int(drawingId),
		s:              s,
		p:              p,
		content:        make([]interface{}, 0, 32),
		completed:      false}
	s.ds = append(s.ds, drawing)
	return drawing, nil
}

func (s *stack) IsComplete() bool {
	return s.completed
}

func (s *stack) Complete() error {
	if s.IsComplete() {
		return errors.New("stack is already completed")
	}
	if !s.validateStack(s) {
		return errors.New("cannot complete stack; is it empty?")
	}

	if err := s.prepStatement(
		"completeStack",
		`UPDATE stacks
    SET complete = TRUE
    WHERE sid = ?;`,
		&s.stackBackend.completeStack); err != nil {
		return err
	}

	if err := s.prepStatement(
		"wipeOrder",
		`DELETE FROM stackHoldings
    WHERE sid = ?;`,
		&s.stackBackend.wipeOrder); err != nil {
		return err
	}

	conn := s.Conn()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	//TODO: revisit to see if these fields can just be autopromoted
	complete := tx.Stmt(s.stackBackend.completeStack)
	wipe := tx.Stmt(s.stackBackend.wipeOrder)
	if _, err := complete.Exec(s.sid); err != nil {
		return err
		tx.Rollback()
	}

	if _, err := wipe.Exec(s.sid); err != nil {
		return err
		tx.Rollback()
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	s.completed = true
	return nil
}
