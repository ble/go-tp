package model

type Game interface {
	Gid() string
	Players() []Player
	NextPlayer(Player) Player
	JoinGame(User, string) (Player, error)
	PlayerForId(string) Player

	Stacks() []Stack
	StacksInProgress() map[Player][]Stack
	StacksFor(Player) []Stack
	PassStack(Player) (Stack, error)

	IsStarted() bool
	Start() error
	IsComplete() bool
	Complete() error
}
