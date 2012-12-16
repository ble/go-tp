package persistence

import (
	"ble/tpg/model"
	"fmt"
)

type playerBackend struct {
	*Backend
}

type player struct {
	*playerBackend
	user      model.User
	pseudonym string
	pid       string
	game      model.Game
}

func (p *playerBackend) DumpAllPlayers() {
	stmt, err := p.Conn().Prepare(
		`SELECT pid, pseudonym, gid, uid, playOrder
        FROM players`)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := stmt.Query()

	if err != nil {
		fmt.Println("e2")
		return
	}
	defer rows.Close()

	var pid, pseudonym, gid, uid string
	var playOrder int
	for rows.Next() {
		err = rows.Scan(&pid, &pseudonym, &gid, &uid, &playOrder)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(pid, pseudonym, gid, uid, playOrder)
	}
}

func (p *player) User() model.User {
	return p.user
}

func (p *player) Pseudonym() string {
	return p.pseudonym
}

func (p *player) Pid() string {
	return p.pid
}
