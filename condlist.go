package queue

import (
	"container/list"
	"sync"
)

// An implementation of Queue using a double linked list protected by
// a condition variable.
type CondListQueue struct {
	// Condition Variable protecting access this `q`.
	cond *sync.Cond

	// Elements in the queue.
	q *list.List
}

func NewCondListQueue() *CondListQueue {
	return &CondListQueue{
		cond: &sync.Cond{L: &sync.Mutex{}},
		q: list.New(),
	}
}

func (this *CondListQueue) PopFront() interface{} {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	for this.q.Len() == 0 {
		this.cond.Wait()
	}
	e := this.q.Front()
	defer this.q.Remove(e)
	return e.Value
}

func (this *CondListQueue) PushBack(v interface{}) {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	this.q.PushBack(v)
	this.cond.Signal()
}
