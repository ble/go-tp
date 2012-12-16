package model

type Drawings interface {
	DrawingById(string)
	CanDrawingBeSeen(Drawing, User)
}
