package persistence

import (
	"errors"
)

func (b *Backend) createTables() error {
	var tableCreationStatements []string = []string{
		`CREATE TABLE users (
      uid    INTEGER PRIMARY KEY,
      email  TEXT(255) NOT NULL CONSTRAINT uniqueEmail UNIQUE ON CONFLICT FAIL,
      alias  TEXT(32)  NOT NULL CONSTRAINT uniqueAlias UNIQUE ON CONFLICT FAIL,
      pwHash TEXT(20)  NOT NULL);`,

		`CREATE INDEX userByAlias ON users (alias)`,

		`CREATE TABLE games (
      gid       STRING PRIMARY KEY,
      started   BOOLEAN
      completed BOOLEAN
      roomName  TEXT(255) NOT NULL);`,

		`CREATE TABLE players (
      pid        INTEGER PRIMARY KEY,
      pseudonym  TEXT(64) NOT NULL,
      gid        INTEGER REFERENCES games (gid),
      uid        INTEGER REFERENCES users (uid),
      playOrder  INTEGER,
      CONSTRAINT uniqueNamePerGame UNIQUE (gid, pseudonym) ON CONFLICT FAIL,
      CONSTRAINT uniqueOrder UNIQUE (playOrder) ON CONFLICT FAIL);`,

		`CREATE TABLE stacks (
      sid        INTEGER PRIMARY KEY,
      gid        INTEGER REFERENCES games (gid),
      complete   BOOLEAN,
      holdingPid INTEGER REFERENCES players (pid));`,

		`CREATE TABLE stackHoldings (
      pid        INTEGER REFERENCES players (pid),
      gid        INTEGER REFERENCES games   (gid),
      sid        INTEGER REFERENCES stacks  (sid),
      ord        INTEGER,
      PRIMARY KEY (sid),
      CONSTRAINT uniqueStackOrder UNIQUE (pid, ord) on CONFLICT FAIL);`,

		`CREATE TABLE drawings (
      did          INTEGER PRIMARY KEY,
      sid          INTEGER REFERENCES stacks (sid),
      pid          INTEGER REFERENCES players (pid),
      stackOrder   INTEGER NOT NULL,
      completeJson BLOB,
      complete     BOOLEAN,
      CONSTRAINT uniqueOrderStack UNIQUE (sid, stackOrder) ON CONFLICT FAIL);`,

		`CREATE TABLE drawParts (
      did         INTEGER REFERENCES drawings (did),
      ord         INTEGER,
      json        BLOB,
      PRIMARY KEY (did, ord));`,
	}
	for _, qString := range tableCreationStatements {
		_, err := b.conn.Exec(qString)
		if err != nil {
			return errors.New(err.Error() + "\nSQL = " + qString)
		}
	}
	return nil

}
