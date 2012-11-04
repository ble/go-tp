package drawing

import (
	. "testing"
)

func TestAgent(t *T) {
	a := NewDrawingHandle()
	e0 := a.Draw(DefaultDrawPart)
	b, e1 := a.Read()
	e2 := a.Close()
	t.Log(a, e0, string(b), e1, e2)
}
