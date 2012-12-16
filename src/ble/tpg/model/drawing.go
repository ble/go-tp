package model

type Drawing interface {
	Did() string
	Stack() Stack
	Player() Player
	Content() []interface{}

	Add(interface{}) error

	IsComplete() bool
	Complete() error
}
