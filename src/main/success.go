package main

import (
	"errors"
	"fmt"
	"time"
)

type success chan bool

func (s success) succeededIn(timeout time.Duration) (bool, error) {
	fmt.Println(timeout)
	timeoutCh := time.After(timeout)
	select {
	case <-timeoutCh:
		fmt.Println("a")
		return false, errors.New("timeout")
	case b := <-s:
		fmt.Println("b")
		return b, nil
	}
	fmt.Println("c")
	return false, errors.New("fallthrough")
}
