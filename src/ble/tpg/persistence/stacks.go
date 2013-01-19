package persistence

import (
	"ble/tpg/model"
	"sync"
)

type stacksBackend struct {
	*Backend
	stackById map[string]model.Stack
	*sync.RWMutex
}

func newStacksBackend(b *Backend) *stacksBackend {
	return &stacksBackend{b, make(map[string]model.Stack), &sync.RWMutex{}}
}

func (s *stacksBackend) GetStackForId(sid string) (model.Stack, bool) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	stack, ok := s.stackById[sid]
	return stack, ok
}

func (s *stacksBackend) recordNewStack(stack model.Stack) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.stackById[stack.Sid()] = stack
}
