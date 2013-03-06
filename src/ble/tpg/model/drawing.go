package model

import "ble/tpg/drawing"
import "encoding/json"

type Drawing interface {
	Did() string
	Stack() Stack
	Player() Player
	Content() []json.Marshaler

	Add(drawing.DrawPart) error

	IsComplete() bool
	Complete() error
}
