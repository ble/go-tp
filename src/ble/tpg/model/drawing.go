package model

type Drawing interface {
	Stack() Stack
	Player() Player
	Content() []interface{}

	Add(interface{}) error

	IsComplete() bool
	Complete() error
}
