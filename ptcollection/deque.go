package ptcollection

import (
	"GoStudy/ptlog"
	"errors"
)

type Deque struct {
	head      int
	tail      int
	size      int
	container []interface{}
}

func NewDeque(args ...interface{}) Deque {
	return Deque{
		head:      0,
		tail:      len(args) - 1,
		size:      len(args),
		container: args,
	}
}

func (deque Deque)LPush(value interface{}) interface{} {
	if deque.size==0{
		ptlog.Error(errors.New("deque is null"))
		return nil
	}
	deque.tail--
	temp:=deque.container[0]
	copy(deque.container,deque.container[1:])
	deque.size--
	return temp
}
