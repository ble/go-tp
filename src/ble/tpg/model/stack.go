package model

type Stack interface {
	Sid() int
	Game() Game
	AllDrawings() []Drawing

	AddDrawing(Player) (Drawing, error)

	IsComplete() bool
	Complete() error
}
