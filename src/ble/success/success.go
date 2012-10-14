package success

import (
	"errors"
	"time"
)

type Success chan bool

var Second time.Duration = time.Second

func (s Success) SucceededIn(timeout time.Duration) error {
	timeoutCh := time.After(timeout)
	select {
	case <-timeoutCh:
		return errors.New("timeout")
	case _ = <-s:
		return nil
	}
	return errors.New("fallthrough")
}
