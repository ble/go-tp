package persistence

import (
	"ble/hash"
	"ble/tpg/model"
	"database/sql"
	"sync"
	"testing"
)

var whatIsAMan string = "A mIsErAbLe PiLe Of SeCreTs"

type Backend struct {
	conn    *sql.DB
	loggers []*testing.T

	*drawingBackend
	*stackBackend
	*gameBackend
	*userBackend
	*playerBackend
	*gamesBackend
	*drawingsBackend
}

func (b *Backend) Conn() *sql.DB {
	return b.conn
}

func OpenBackend(filename string) (*Backend, error) {
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	b := &Backend{conn: conn, loggers: []*testing.T{}}
	b.drawingBackend = &drawingBackend{Backend: b}
	b.stackBackend = &stackBackend{Backend: b}
	b.gameBackend = &gameBackend{Backend: b}
	b.userBackend = &userBackend{Backend: b}
	b.playerBackend = &playerBackend{Backend: b}
	b.gamesBackend = &gamesBackend{
		Backend:  b,
		allGames: make(map[string]model.Game),
		RWMutex:  new(sync.RWMutex)}
	b.drawingsBackend = &drawingsBackend{
		Backend:  b,
		drawings: make(map[string]model.Drawing),
		RWMutex:  new(sync.RWMutex)}
	return b, nil
}

func NewBackend(filename string) (*Backend, error) {
	b, err := OpenBackend(filename)
	if err != nil {
		return nil, err
	}
	err = b.createTables()
	if err != nil {
		b.Conn().Close()
		return nil, err
	}
	return b, nil
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

func (b Backend) hashPw(alias, pw string) string {
	h := hash.NewHashEasy()
	return h.WriteStrAnd(pw).
		WriteStrAnd(alias).
		WriteStrAnd(whatIsAMan).String()
}
