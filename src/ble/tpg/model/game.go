package model

type Game interface {
	Players() []Player
	NextPlayer(Player) Player
	JoinGame(User, string) (Player, error)

	Stacks() []Stack
	StacksInProgress() map[Player][]Stack

	IsStarted() bool
	Start() error
	IsComplete() bool
	Complete() error
}
