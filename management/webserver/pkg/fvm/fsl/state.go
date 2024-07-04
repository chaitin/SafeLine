package fsl

import (
	"chaitin.cn/dev/go/errors"
)

type State int

const (
	S_NOP State = iota
	S_ABORT
	S_RETURN
)

func (s State) String() string {
	switch s {
	case S_NOP:
		return "NOP"
	case S_ABORT:
		return "ABORT"
	case S_RETURN:
		return "RETURN"
	default:
		panic(errors.New("wrong State"))
	}
}
