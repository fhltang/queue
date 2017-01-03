package queue

// An interface for a thread-safe FIFO Queue.
type Queue interface {
	// Insert into the queue.  Should never block; conceptually
	// the queue has unbounded capacity.
	PushBack(interface{})

	// Remove from the queue.
	//
	// Blocks if the queue is empty.
	PopFront() interface{}
}

