package model

type Player interface {
	User() User
	Pseudonym() string
}
