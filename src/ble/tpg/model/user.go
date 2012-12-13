package model

type User interface {
	Alias() string
	Email() string
	Uid() int
}
