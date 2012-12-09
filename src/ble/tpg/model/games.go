package model

type Games interface {
	AllGames() []Game
	CreateGame(roomName string) Game
}

type Stacks interface {
	CreateStack(Game, Player)
}

type Player interface {
	GetUser() User
	Pseudonym() string
	Stacks() []Stack
}

type Stack interface {
	TopDrawing() Drawing
	DrawingCount() int
	AllDrawings() []Drawing
	IsComplete() bool

	AddDrawing() Drawing
	Complete()
}

type Drawing interface {
	Add(interface{})
}

type Game interface {
	RoomName() string
	Players() []Player
	Stacks() map[Player]Stack

	PassStack(Player)
	JoinGame(User) Player
	Complete()
}
