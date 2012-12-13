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
