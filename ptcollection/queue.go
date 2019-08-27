package ptcollection

type queueUnit struct {
	next  *queueUnit
	value interface{}
}
type queue struct {
	head *queueUnit
	tail *queueUnit
	size int
}

//NewLinkedQueue 基于链表实现queue，2000100（1万次插入后1万次取出的时间）
func NewLinkedQueue() *queue {
	return &queue{size: 0}
}

func (q *queue) Push(arg interface{}) {
	unit := queueUnit{value: arg}

	if q.size == 0 {
		q.head = &unit
		q.tail = &unit
	} else {
		q.tail.next = &unit
		q.tail = &unit
	}
	q.size++
}
func (q *queue) Pop() interface{} {
	if q.size == 0 {
		return nil
	}
	val := q.head.value
	q.head = q.head.next
	q.size--
	return val
}
func (q *queue) IsEmpty() bool {
	return q.size == 0
}

type arrayQueue struct {
	content []interface{}
}

//NewArrayQueue 使用数组实现队列 1000100（1万次插入后1万次取出的时间）
func NewArrayQueue() *arrayQueue {
	return &arrayQueue{}
}
func (aq *arrayQueue) Push(arg interface{}) {
	aq.content = append(aq.content, arg)
}
func (aq *arrayQueue) Pop() interface{} {
	if len(aq.content) == 0 {
		return nil
	}
	val := aq.content[0]
	aq.content = aq.content[1:]
	return val
}
func (aq *arrayQueue) IsEmpty() bool {
	return len(aq.content)==0
}
