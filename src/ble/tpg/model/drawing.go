package model

import "encoding/json"

type Drawing interface {
	Stack() Stack
	Player() Player
	Content() json.Marshaler

	Add(json.Marshaler) error

	IsComplete() bool
	Complete() error
}
