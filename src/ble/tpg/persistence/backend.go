package persistence

import (
	"ble/hash"
	"database/sql"
	"fmt"
	"testing"
)

var whatIsAMan string = "A mIsErAbLe PiLe Of SeCreTs"

type Backend struct {
	conn                                                                                             *sql.DB
	loggers                                                                                          []*testing.T
	countPlayersInGame, createPlayer, createGame, createUser, getUserByAlias, logInUser, getAllGames *sql.Stmt
}

func NewBackend(filename string) (*Backend, error) {
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	b := Backend{conn: conn, loggers: []*testing.T{}}
	return &b, nil
}

func (b *Backend) prepAllStatements() error {
	type statementInternal struct {
		desc string
		stmt **sql.Stmt
		cmd  string
	}
	statements := []statementInternal{
		{"createUser",
			&b.createUser,
			`INSERT INTO users (email, alias, pwHash) 
       VALUES (?, ?, ?)`},
		{"getUserByAlias",
			&b.getUserByAlias,
			`SELECT * FROM users WHERE alias == ?`},
		{"logInUser",
			&b.logInUser,
			`SELECT uid, email FROM users
      WHERE alias = ? and pwHash = ?;`},
		{"countPlayersInGame",
			&b.countPlayersInGame,
			`SELECT COUNT(pid) FROM players WHERE gid = ?;`},
		{"createPlayer",
			&b.createPlayer,
			`INSERT INTO players (pseudonym, playOrder, uid, gid)
       VALUES (?, ?, ?, ?)`}}
	for i := range statements {
		s := statements[i]
		fmt.Println(s)
		fmt.Println(*s.stmt)
		err := b.prepStatement(s.desc, s.cmd, s.stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Backend) prepStatement(desc, sql string, stmt **sql.Stmt) error {
	if *stmt == nil {
		newStmt, err := b.conn.Prepare(sql)
		if err != nil {
			b.logError("Preparing statment `"+desc+"`", err)
			return err
		}
		*stmt = newStmt
	}
	return nil
}

func (b Backend) validateRoomName(roomName string) error {
	return nil
}

func (b Backend) validateEmail(email string) error {
	return nil
}

func (b Backend) validateAlias(alias string) error {
	return nil
}

func (b Backend) validatePassword(pw string) error {
	return nil
}

func (b *Backend) RegisterLogger(t *testing.T) {
	b.loggers = append(b.loggers, t)
}

func (b Backend) logError(loc string, args ...interface{}) {
	allArgs := append([]interface{}{loc}, args)
	for _, l := range b.loggers {
		l.Log(allArgs...)
	}
}

func (b Backend) hashPw(pw string) string {
	h := hash.NewHashEasy()
	return h.WriteStrAnd(pw).WriteStrAnd(whatIsAMan).String()
}
