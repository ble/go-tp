package ephemeral

import "fmt"

type IdGenerator interface {
	NewId() string
}

type trivialIdGenerator struct {
	count uint64
}

func (t *trivialIdGenerator) NewId() string {
	c0 := t.count
	t.count++
	return fmt.Sprintf("%d", c0)
}
