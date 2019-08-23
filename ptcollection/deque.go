package ptcollection

import (
	"errors"

	//pt go lib
	"github.com/chenhnu/go-lib/ptlog"
)

type deque struct {
	head      *interface{}
	tail      *interface{}
	size      int
	container []interface{}
}

func NewDeque(args ...interface{}) *deque {
	if args == nil {
		return &deque{}
	}
	return &deque{
		head:      &args[0],
		tail:      &args[len(args)-1],
		size:      len(args),
		container: args,
	}
}
func (deq *deque) GetContent() []interface{} {
	return deq.container
}
func (deq *deque) IsEmpty() bool {
	if deq.size == 0 {
		return true
	}
	return false
}
func (deq *deque) Size() int {
	return deq.size
}
func (deq *deque) LPop() interface{} {
	if deq.size == 0 {
		ptlog.Error(errors.New("deque is null"))
		return nil
	}
	res := &deq.head
	deq.container = deq.container[1:]
	deq.size--
	deq.head = &deq.container[0]
	return res
}
func (deq *deque) RPop() interface{} {
	if deq.size == 0 {
		ptlog.Error(errors.New("deque is null"))
		return nil
	}
	res := &deq.tail
	deq.container = deq.container[0 : deq.size-1]
	deq.size--
	deq.tail = &deq.container[deq.size-1]
	return res
}

func (deq *deque) LPush(value interface{}) {
	deq.head = &value
	deq.size++
	tem := []interface{}{value}
	deq.container = append(tem, deq.container...)
}
func (deq *deque) RPush(value interface{}) {
	deq.tail = &value
	deq.size++
	deq.container = append(deq.container, value)
}
