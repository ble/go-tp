package model

type Games interface {
	AllGames() map[string]Game
	CreateGame(string) (Game, error)
}
