package queue

import (
	"sync"
)

// An implementation of Queue using two slices (used as stacks)
// protected by a condition variable.
type CondSliceQueue struct {
	// Condition Variable protecting the two slices.
	cond *sync.Cond

	// Elements at the front of queue in reverse order.
	revFront []interface{}

	// Elements at the back fo the queue.
	back []interface{}
}

func NewCondSliceQueue() *CondSliceQueue {
	return &CondSliceQueue{cond: &sync.Cond{L: &sync.Mutex{}}}
}

func (this *CondSliceQueue) PopFront() interface{} {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()

	for len(this.revFront) == 0 && len(this.back) == 0 {
		this.cond.Wait()
	}

	if len(this.revFront) == 0 {
		// Reverse `this.back` in place and swap with `this.revFront`.
		for l, r := 0, len(this.back)-1; l < r; l, r = l+1, r-1 {
			this.back[l], this.back[r] = this.back[r], this.back[l]
		}
		this.back, this.revFront = this.revFront, this.back		
	}

	v := this.revFront[len(this.revFront)-1]
	this.revFront = this.revFront[:len(this.revFront)-1]

	return v
}

func (this *CondSliceQueue) PushBack(v interface{}) {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	this.back = append(this.back, v)
	this.cond.Signal()
}
