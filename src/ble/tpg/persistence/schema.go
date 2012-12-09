package persistence

import (
	. "database/sql"
)

func createTables(b Backend) ([]Result, []error) {
	rs := make([]Result, len(tableCreationStatements))
	es := make([]error, len(tableCreationStatements))
	for ix, qString := range tableCreationStatements {
		rs[ix], es[ix] = b.conn.Exec(qString)
	}
	return rs, es
}

var tableCreationStatements []string = []string{
	`CREATE TABLE users (
    uid    INTEGER PRIMARY KEY,
    email  TEXT(255) NOT NULL CONSTRAINT uniqueEmail UNIQUE ON CONFLICT FAIL,
    alias  TEXT(32)  NOT NULL CONSTRAINT uniqueAlias UNIQUE ON CONFLICT FAIL,
    pwHash TEXT(20)  NOT NULL);`,

	`CREATE INDEX userByAlias ON users (alias)`,

	`CREATE TABLE games (
    gid      INTEGER PRIMARY KEY,
    roomName TEXT(255) NOT NULL);`,

	`CREATE TABLE players (
    pid        INTEGER PRIMARY KEY,
    pseudonym  TEXT(64) NOT NULL,
    gid        INTEGER REFERENCES games (gid),
    uid        INTEGER REFERENCES users (uid),
    CONSTRAINT uniqueNamePerGame UNIQUE (gid, pseudonym) ON CONFLICT FAIL);`,

	`CREATE TABLE gamePlayerOrder (
    gid INTEGER REFERENCES games (gid),
    pid INTEGER REFERENCES players (pid),
    orderInGame INTEGER NOT NULL,
    PRIMARY KEY (gid, orderInGame));`,

	`CREATE TABLE stacks (
    sid        INTEGER PRIMARY KEY,
    gid        INTEGER REFERENCES games (gid),
    holdingPid INTEGER REFERENCES players (pid));`,

	`CREATE TABLE drawings (
    did          INTEGER PRIMARY KEY,
    sid          INTEGER REFERENCES stacks (sid),
    pid          INTEGER REFERENCES players (pid),
    stackOrder   INTEGER NOT NULL,
    completeJson BLOB,
    CONSTRAINT uniqueOrderInStack UNIQUE (sid, stackOrder) ON CONFLICT FAIL);`}
