package model

type Stack interface {
	Sid() string
	Game() Game
	AllDrawings() []Drawing
	TopDrawing() Drawing

	AddDrawing(Player) (Drawing, error)

	IsComplete() bool
	Complete() error
}
