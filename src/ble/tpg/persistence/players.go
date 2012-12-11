package persistence

import (
	"ble/tpg/model"
)

type player struct {
	model.User
	pseudonym string
}

func (p player) GetUser() model.User {
	return p.User
}

func (p player) Pseudonym() string {
	return p.pseudonym
}
