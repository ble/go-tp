package model

type Games interface {
	AllGames() (map[string]Game, error)
	CreateGame(roomName string) (Game, error)
}

type Stacks interface {
	CreateStack(Game, Player)
}

type Player interface {
	GetUser() User
	Pseudonym() string
}

type Stack interface {
	TopDrawing() Drawing
	DrawingCount() int
	AllDrawings() []Drawing
	IsComplete() bool
	AddDrawing(Player) Drawing
	Complete()
}

type Drawing interface {
	Add(interface{})
}

type Game interface {
	RoomName() string
	Players() []Player
	Stacks() map[Player][]Stack

	PassStack(Player)
	JoinGame(User, string) Player
	Complete()
}
