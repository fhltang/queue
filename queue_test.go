package queue_test

import (
	"fmt"
	"github.com/fhltang/queue"
	"sync"
	"testing"
)

const Bound int = 10000

type TestCase struct {
	Name string
	New func () queue.Queue
}

var TestCases = []TestCase{
	TestCase{"CondListQueue", func() queue.Queue { return queue.NewCondListQueue() }},
	TestCase{"CondSliceQueue", func() queue.Queue { return queue.NewCondSliceQueue() }} ,
	TestCase{"ChanSliceQueue", func() queue.Queue { return queue.NewChanSliceQueue() }},
	TestCase{"BoundedChanQueue", func() queue.Queue { return queue.NewBoundedChanQueue(Bound) }},
}

func TestPushPop(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			q := tc.New()

			q.PushBack(1)

			if item := q.PopFront().(int); 1 != item {
				t.Errorf("Got %d", item)
			}
		})
	}
}

func TestPushPushPopPop(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			q := tc.New()

			q.PushBack(1)
			q.PushBack(2)

			if item := q.PopFront().(int); 1 != item {
				t.Errorf("Got %d", item)
			}
			if item := q.PopFront().(int); 2 != item {
				t.Errorf("Got %d", item)
			}
		})
	}
}

func TestPush3Pop2Push2Pop3(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			q := tc.New()

			q.PushBack(1)
			q.PushBack(2)
			q.PushBack(3)

			if item := q.PopFront().(int); 1 != item {
				t.Errorf("Got %d", item)
			}
			if item := q.PopFront().(int); 2 != item {
				t.Errorf("Got %d", item)
			}

			q.PushBack(4)
			q.PushBack(5)

			if item := q.PopFront().(int); 3 != item {
				t.Errorf("Got %d", item)
			}
			if item := q.PopFront().(int); 4 != item {
				t.Errorf("Got %d", item)
			}
			if item := q.PopFront().(int); 5 != item {
				t.Errorf("Got %d", item)
			}
		})
	}
}

// Benchmark: serially push then serially pop.
func SerialPushSerialPop(q queue.Queue) {
	N := 10000
	for j := 0; j < N; j++ {
		q.PushBack(j)
	}
	for j := 0; j < N; j++ {
		q.PopFront()
	}
}

// Benchmark: push items concurrently and read them serially.
func ConcurrentPushSerialPop(q queue.Queue) {
	N := 10000
	go func() {
		for j := 0; j < N; j++ {
			q.PushBack(j)
		}
	}()
	for j := 0; j < N; j++ {
		q.PopFront()
	}
}

// Benchmark: push items serially and read them concurrently.
func SerialPushConcurrentPop(q queue.Queue) {
	N := 10000
	w := sync.WaitGroup{}
	w.Add(N)
	for j := 0; j < N; j++ {
		go func() {
			q.PopFront()
			w.Done()
		}()
	}
	for j := 0; j < N; j++ {
		q.PushBack(j)
	}
	w.Wait()
}

// Benchmark: alternately push and pop items in batches of 100.
func AlternatePushPop(q queue.Queue) {
	M := 100
	N := 100
	for i := 0; i < N; i ++{
		for j := 0; j < M; j++ {
			q.PushBack(i*M + j)
		}
		for j := 0; j < M; j++ {
			q.PopFront()
		}
	}
}

type Benchmark struct {
	Name string
	Run func (queue.Queue)
}

var Benchmarks = []Benchmark{
	Benchmark{"New", func(q queue.Queue) { return } },
	Benchmark{"SerialPushSerialPop", SerialPushSerialPop},
	Benchmark{"ConcurrentPushSerialPop", ConcurrentPushSerialPop},
	Benchmark{"SerialPushConcurrentPop", SerialPushConcurrentPop},
	Benchmark{"AlternatePushPop", AlternatePushPop},
}

func BenchmarkAll(b *testing.B) {
	for _, bm := range Benchmarks {
		for _, tc := range TestCases {
			b.Run(fmt.Sprintf("%s_%s", bm.Name, tc.Name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bm.Run(tc.New())
				}
			})
		}
	}
}

