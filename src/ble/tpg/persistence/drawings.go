package persistence

import (
	"ble/tpg/model"
	"sync"
)

type drawingsBackend struct {
	*Backend
	drawings map[string]model.Drawing
	*sync.RWMutex
}

func (d *drawingsBackend) addDrawingToBackend(drawing model.Drawing) {
	d.Lock()
	d.drawings[drawing.Did()] = drawing
	defer d.Unlock()
}

func (d *drawingsBackend) GetDrawingForId(did string) (model.Drawing, bool) {
	d.RLock()
	drawing, present := d.drawings[did]
	d.RUnlock()
	if !present {
		return nil, false
	}
	return drawing, true
}

func (dbe *drawingsBackend) CanDrawingBeSeen(d model.Drawing, u model.User) bool {
	if !d.IsComplete() {
		return u.Uid() == d.Player().User().Uid()
	}
	stack := d.Stack()
	if stack.IsComplete() {
		return true
	}
	for _, drawing := range stack.AllDrawings() {
		user := drawing.Player().User()
		if user.Uid() == u.Uid() {
			return true
		}
	}
	return false
}
