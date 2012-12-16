package persistence

import (
	"errors"
)

func (b *Backend) createTables() error {
	var tableCreationStatements []string = []string{
		`CREATE TABLE users (
      uid    TEXT(24) PRIMARY KEY NOT NULL,
      email  TEXT(255) NOT NULL CONSTRAINT uniqueEmail UNIQUE ON CONFLICT FAIL,
      alias  TEXT(32)  NOT NULL CONSTRAINT uniqueAlias UNIQUE ON CONFLICT FAIL,
      pwHash TEXT(20)  NOT NULL);`,

		`CREATE INDEX userByAlias ON users (alias)`,

		`CREATE TABLE games (
      gid       TEXT(24) PRIMARY KEY NOT NULL,
      started   BOOLEAN,
      complete  BOOLEAN,
      roomName  TEXT(255) NOT NULL);`,

		`CREATE TABLE players (
      pid        TEXT(24) PRIMARY KEY NOT NULL,
      pseudonym  TEXT(64) NOT NULL,
      gid        TEXT(24) REFERENCES games (gid),
      uid        TEXT(24) REFERENCES users (uid),
      playOrder  INTEGER,
      CONSTRAINT uniqueNamePerGame   UNIQUE (gid, pseudonym) ON CONFLICT FAIL,
      CONSTRAINT uniquePlayerPerUser UNIQUE (gid, uid) ON CONFLICT FAIL,
      CONSTRAINT uniqueOrder         UNIQUE (gid, playOrder) ON CONFLICT FAIL);`,

		`CREATE TABLE stacks (
      sid        TEXT(24) PRIMARY KEY NOT NULL,
      gid        TEXT(24) REFERENCES games (gid),
      complete   BOOLEAN,
      holdingPid TEXT(24) REFERENCES players (pid));`,

		`CREATE TABLE stackHoldings (
      pid        TEXT(24) REFERENCES players (pid),
      gid        TEXT(24) REFERENCES games   (gid),
      sid        TEXT(24) REFERENCES stacks  (sid),
      ord        INTEGER,
      PRIMARY KEY (sid),
      CONSTRAINT uniqueStackOrder UNIQUE (pid, ord) on CONFLICT FAIL);`,

		`CREATE TABLE drawings (
      did          TEXT(24) PRIMARY KEY NOT NULL,
      sid          TEXT(24) REFERENCES stacks (sid),
      pid          TEXT(24) REFERENCES players (pid),
      stackOrder   INTEGER NOT NULL,
      completeJson BLOB,
      complete     BOOLEAN,
      CONSTRAINT uniqueOrderStack UNIQUE (sid, stackOrder) ON CONFLICT FAIL);`,

		`CREATE TABLE drawParts (
      did         TEXT(24) REFERENCES drawings (did),
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
