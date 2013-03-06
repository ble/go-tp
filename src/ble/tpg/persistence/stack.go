package persistence

import (
	"ble/hash"
	"ble/tpg/model"
	"database/sql"
	"encoding/json"
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
	sid       string
	g         model.Game
	ds        []model.Drawing
	completed bool
}

func typeCheckStack() model.Stack {
	return &stack{}
}

func (s *stack) Sid() string {
	return s.sid
}

func (s *stack) Game() model.Game {
	return s.g
}

func (s *stack) AllDrawings() []model.Drawing {
	return s.ds
}

func (s *stack) TopDrawing() model.Drawing {
	return s.ds[len(s.ds)-1]
}

func (s *stack) AddDrawing(p model.Player) (model.Drawing, error) {
	if err := s.prepStatement(
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
		&s.addDrawing); err != nil {
		return nil, err
	}
	drawingId := hash.EasyNonce(p.Pid(), s.Sid())
	_, err := s.addDrawing.Exec(
		drawingId,
		s.sid,
		p.Pid(),
		len(s.AllDrawings()),
		s.sid)
	if err != nil {
		return nil, err
	}
	drawing := &drawing{
		drawingBackend: s.drawingBackend,
		did:            drawingId,
		s:              s,
		p:              p,
		content:        make([]json.Marshaler, 0, 32),
		completed:      false}
	s.ds = append(s.ds, drawing)
	s.addDrawingToBackend(drawing)
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
