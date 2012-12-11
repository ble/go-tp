package persistence

/*
import (
	"ble/tpg/model"
  "ble/hash"
)

type games struct {
  allGames map[string]model.Game
  requests chan<- interface{}
  Backend
}

type stack struct {}

type game struct {
  gid, roomName string
  players []model.Player
  stacks map[model.Player][]model.Stack
  isComplete bool
}

func (g *game) RoomName() string {
  return g.roomName
}

func (g *game) Players() []model.Player {
  return g.players
}

func (g *game) Stacks() map[model.Player][]model.Stack {
  return g.stacks
}

func (g *game) PassStack(model.Player) {
}

func (g *game) JoinGame(u model.User, pseudonym string) (model.Player, ) {
  if err := g.Backend.prepStatement(
    "countPlayersInGame",
    `SELECT COUNT(pid) FROM players WHERE gid == ?;`,
    &g.Backend.countPlayersInGame); err != nil {

  }
  err := g.Backend.prepStatement(
    "createPlayer",

}

func (g *game) Complete() {
  g.isComplete = true
}


func (g *games) CreateGame(roomName string) (model.Game, error) {
  if err := g.Backend.prepStatement(
    "createGame",
    `INSERT INTO games (gid, roomName) VALUES (?, ?);`,
    &g.Backend.createGame); err != nil {
    return nil, err
  }
  if err := g.Backend.validateRoomName(roomName); err != nil {
    return nil, err
  }
  gid := hash.NewHashEasy().Nonce().String()
  if _, err := g.Backend.createGame.Exec(gid, roomName); err != nil {
    return nil, err
  }
  newGame := game{gid: gid, roomName: roomName}
  g.allGames[gid] = &newGame
  return &newGame, nil
}

/*
type games struct {
	allGames map[string]game
	requests chan<- interface{}
}

type game struct {
	gid, roomName string
	players       map[int]player
	stacks        []interface{}
	Backend
}

func (g *game) RoomName() string {
	return g.roomName
}

func (g *game) getPlayers() {
  err := g.Backend.prepStatement(
    "getPlayers",
    "SELECT players.pid, players.pseudonym, gamePlayerOrder
}

func (g *game) Stacks() map[Player][]Stack {
	g.getStacks()
	if g.stacks == nil {
	}
	return nil
}

func (g *game) getStacks() {
	g.getPlayers()
	err := g.Backend.prepStatement(
		"getGameStacks",
		`SELECT stacks.sid, stacks.holdingPid, drawings.did, drawings.stackOrder
     FROM stacks JOIN drawings, stacks USING sid
     WHERE gid == ?`,
		&g.Backend.getGameStacks)
	if g.stacks == nil {
		rows, err := g.Backend.getGameStacks.Query(g.gid)
		if err != nil {
			g.Backend.logError("getting game stacks", err)
			return
		}
		g.stacks = make(map[player][]stack)
		defer rows.Close()
		for rows.Next() {
		}
	}
}

func (g *game) Players() []Player {

}

func (b *Backend) AllGames() ([]model.Game, error) {
	if err := b.prepStatement(
		"getAllGames",
		`SELECT gid, roomName from games`,
		&b.getAllGames); err != nil {
		return nil, err
	}
	rows, err := b.getAllGames.Query()
	if err != nil {
		return nil, err
	}
	games := make([]game, 0, 128)
}
*/
