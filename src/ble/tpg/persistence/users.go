package persistence

import (
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type user struct {
	uid                  int
	email, alias, pwHash string
	b                    *Backend
}

func (u user) Alias() string {
	return u.alias
}

func (u user) Email() string {
	return u.email
}

func (u user) Uid() int {
	return u.uid
}

func (b *Backend) LogInUser(alias, pw string) (model.User, error) {
	pwHash := b.hashPw(alias, pw)
	b.prepStatement(
		"logInUser",
		`SELECT uid, email FROM users
     WHERE alias = ? AND pwHash = ?;`,
		&b.logInUser)
	row := b.logInUser.QueryRow(alias, pwHash)
	var uid int
	var email string
	err := row.Scan(&uid, &email)
	if err == sql.ErrNoRows {
		return nil, errors.New("bad alias or password")
	}
	return user{uid, email, alias, pwHash, b}, nil
}

func (b *Backend) CreateUser(email, alias, pw string) (model.User, error) {
	err := b.validateEmail(email)
	if err != nil {
		return nil, err
	}

	err = b.validateAlias(alias)
	if err != nil {
		return nil, err
	}

	err = b.validatePassword(pw)
	if err != nil {
		return nil, err
	}

	pwHash := b.hashPw(alias, pw)
	if err = b.prepStatement(
		"createUser",
		`INSERT INTO users
         (email, alias, pwHash) VALUES (?, ?, ?)`,
		&b.createUser); err != nil {
		return nil, err
	}

	if err = b.prepStatement(
		"getUserByAlias",
		`SELECT * FROM users WHERE alias == ?`,
		&b.getUserByAlias); err != nil {
		return nil, err
	}
	tx, err := b.conn.Begin()
	if err != nil {
		return nil, err
	}
	insert := tx.Stmt(b.createUser)
	read := tx.Stmt(b.getUserByAlias)
	_, err = insert.Exec(email, alias, pwHash)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	row := read.QueryRow(alias)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var uid int
	err = row.Scan(&uid, &email, &alias, &pwHash)
	if err != nil {
		return nil, err
	}
	return user{uid, email, alias, pwHash, b}, nil
}
