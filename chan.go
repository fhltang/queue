package queue

// An implementation of Queue using a bounded channel.
//
// This is an inaccurate implementation of the contract of Queue since
// it is bounded and PushBack() can block.  However, for a
// sufficiently large bound, it is almost correct.
type BoundedChanQueue struct {
	q chan interface{}
}

func NewBoundedChanQueue(n int) *BoundedChanQueue {
	return &BoundedChanQueue{q: make(chan interface{}, n)}
}

func (this *BoundedChanQueue) PopFront() interface{} {
	return <-this.q
}

func (this *BoundedChanQueue) PushBack(v interface{}) {
	this.q <- v
}

