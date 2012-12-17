package persistence

import (
	"ble/hash"
	"ble/tpg/model"
	"database/sql"
	"errors"
)

type userBackend struct {
	*Backend
	createUser, getUserByAlias, logInUser *sql.Stmt
}

type usersBackend struct {
	*Backend
	allUsers map[string]model.User
}

type user struct {
	*userBackend
	uid, email, alias, pwHash string
	b                         *Backend
}

func (u user) Alias() string {
	return u.alias
}

func (u user) Email() string {
	return u.email
}

func (u user) Uid() string {
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
	var uid, email string
	err := row.Scan(&uid, &email)
	if err == sql.ErrNoRows {
		return nil, errors.New("bad alias or password")
	}
	return user{b.userBackend, uid, email, alias, pwHash, b}, nil
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
	userId := hash.EasyNonce(alias, pwHash, email)

	if err = b.prepStatement(
		"createUser",
		`INSERT INTO users
         (uid, email, alias, pwHash) VALUES (?, ?, ?, ?)`,
		&b.createUser); err != nil {
		return nil, err
	}
	_, err = b.createUser.Exec(userId, email, alias, pwHash)
	if err != nil {
		return nil, err
	}
	newUser := user{b.userBackend, userId, email, alias, pwHash, b}
	b.allUsers[newUser.Uid()] = newUser
	return newUser, nil
}

func (b *Backend) GetUserById(uid string) (model.User, error) {
	user, present := b.allUsers[uid]
	if !present {
		return nil, errors.New("no such user")
	}
	return user, nil
}
