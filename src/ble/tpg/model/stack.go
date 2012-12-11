package model

type Stack interface {
	Game() Game
	AllDrawings() []Drawing

	AddDrawing(Player) (Drawing, error)

	IsComplete() bool
	Complete() error
}
