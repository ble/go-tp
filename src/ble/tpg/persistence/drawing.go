package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"encoding/json"
	"errors"
)

type drawingBackend struct {
	*Backend
	addDrawPart, completeDrawing, wipeDrawParts *sql.Stmt
}

func (d drawingBackend) validateDrawPart(interface{}) bool {
	return true
}

func (d drawingBackend) validateDrawing(model.Drawing) bool {
	return true
}

type drawing struct {
	*drawingBackend
	did       int
	s         model.Stack
	p         model.Player
	content   []interface{}
	completed bool
}

func typeCheckDrawing() model.Drawing {
	return &drawing{}
}

func (d drawing) Stack() model.Stack {
	return d.s
}

func (d drawing) Player() model.Player {
	return d.p
}

func (d drawing) Content() []interface{} {
	return d.content
}

func (d drawing) IsComplete() bool {
	return d.completed
}

func (d *drawing) Complete() error {
	if d.completed {
		return errors.New("drawing already completed")
	}
	if !d.validateDrawing(d) {
		return errors.New("Drawing can't be completed yet; is it empty?")
	}
	if err := d.prepStatement(
		"completeDrawing",
		`UPDATE drawings
      SET completeJson = ?, complete = TRUE
      WHERE did = ?;`,
		&d.completeDrawing); err != nil {
		return err
	}
	if err := d.prepStatement(
		"wipeDrawParts",
		`DELETE FROM drawParts
      WHERE did = ?;`,
		&d.wipeDrawParts); err != nil {
		return err
	}

	contentJson, err := json.Marshal(d.Content())
	if err != nil {
		return err
	}

	conn := d.Conn()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	complete := tx.Stmt(d.completeDrawing)
	wipe := tx.Stmt(d.wipeDrawParts)
	if _, err := complete.Exec(contentJson, d.did); err != nil {
		return err
		tx.Rollback()
	}
	if _, err := wipe.Exec(d.did); err != nil {
		return err
		tx.Rollback()
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	d.completed = true
	return nil
}

func (d *drawing) Add(x interface{}) error {
	if !d.validateDrawPart(x) {
		return errors.New("invalid draw part")
	}
	if err := d.prepStatement(
		"addDrawPart",
		`INSERT INTO drawParts (did, ord, json)
    SELECT ? as did, COUNT(ord) as ord, ? as json
    FROM drawParts
    WHERE did = ?;`,
		&d.addDrawPart); err != nil {
		return err
	}
	json, err := json.Marshal(x)
	if err != nil {
		return err
	}
	if _, err := d.addDrawPart.Exec(
		d.did,
		d.did,
		json); err != nil {
		return err
	}
	d.content = append(d.content, x)
	return nil
}
