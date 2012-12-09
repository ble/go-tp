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

func (b *Backend) LogInUser(alias, pw string) (model.User, error) {
	pwHash := b.hashPw(pw)
	if b.logInUser == nil {
		logInUser, err := b.conn.Prepare(
			`SELECT uid, email FROM users
       WHERE alias = ? AND pwHash = ?;`)
		if err != nil {
			b.logError("preparing `logInUser` statement", err)
			return nil, err
		}
		b.logInUser = logInUser
	}
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

	pwHash := b.hashPw(pw)
	if b.createUser == nil {
		createUser, err := b.conn.Prepare(
			`INSERT INTO users
         (email, alias, pwHash) VALUES (?, ?, ?)`)
		if err != nil {
			b.logError("preparing `createUser` statement", err)
			return nil, err
		}
		b.createUser = createUser
	}
	if b.getUserByAlias == nil {
		getUserByAlias, err := b.conn.Prepare(
			`SELECT * FROM users WHERE alias == ?`)
		if err != nil {
			b.logError("preparing `getUserByAlias` statement", err)
			return nil, err
		}
		b.getUserByAlias = getUserByAlias
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
