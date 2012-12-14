package persistence

import (
	"ble/tpg/model"
)

type player struct {
	user      model.User
	pseudonym string
	pid       int
	game      model.Game
}

func (p *player) User() model.User {
	return p.user
}

func (p *player) Pseudonym() string {
	return p.pseudonym
}

func (p *player) Pid() int {
	return p.pid
}
