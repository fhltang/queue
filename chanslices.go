package queue

import (
	"runtime"
)

// An implementation of Queue using two slices (used as stacks)
// maintained by a single goroutine.
type ChanSliceQueue struct {
	// Channel for reading from the front of the queue.
	F <-chan interface{}

	// Channel for writing to the back of the queue.
	B chan<- interface{}

	// Channel to signal the goroutine to exit.  This is required
	// to prevent "memory leaks" caused by orphaned goroutines.
	exit chan<- int
}

func NewChanSliceQueue() *ChanSliceQueue {
	front := make(chan interface{})
	back := make(chan interface{})
	exit := make(chan int)

	// Front elements of the queue in reverse order.
	revFrontItems := make([]interface{}, 0)
	// Back elements of the queue.
	backItems := make([]interface{}, 0)

	go func() {
		for {
			if len(revFrontItems) == 0 && len(backItems) == 0 {
				// Special case for when the queue is empty.
				select {
				case <-exit:
					return
				case item := <-back:
					revFrontItems = append(revFrontItems, item)
				}
			}

			// Either `revFrontItems` or `backItems` is
			// non-empty.
			if len(revFrontItems) == 0 {
				// Reverse `backItems` in-place and swap with `revFrontItems`.
				for l, r := 0, len(backItems)-1; l < r; l, r = l+1, r-1 {
					backItems[l], backItems[r] = backItems[r], backItems[l]
				}
				backItems, revFrontItems = revFrontItems, backItems
			}

			// `revFrontItems` is non-empty.
			select {
			case <-exit:
				return
			case front <- revFrontItems[len(revFrontItems)-1]:
				revFrontItems = revFrontItems[:len(revFrontItems)-1]
			case item := <-back:
				backItems = append(backItems, item)
			}
		}
	}()

	queue := &ChanSliceQueue{F: front, B: back, exit: exit}
	runtime.SetFinalizer(queue, func(q *ChanSliceQueue) { q.exit <- 1 })
	return queue
}

func (this *ChanSliceQueue) PopFront() interface{} {
	return <-this.F
}

func (this *ChanSliceQueue) PushBack(v interface{}) {
	this.B <- v
}

