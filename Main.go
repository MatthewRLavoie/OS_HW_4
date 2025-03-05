package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	value int
	next  *Node
}

type QueueMutex struct {
	head, tail         *Node
	headLock, tailLock sync.Mutex
}

func (q *QueueMutex) Init() {
	node := &Node{}
	q.head, q.tail = node, node
}

func (q *QueueMutex) Enqueue(value int) {
	node := &Node{value: value}
	q.tailLock.Lock()
	q.tail.next = node
	q.tail = node
	q.tailLock.Unlock()
}

func (q *QueueMutex) Dequeue(value *int) {
	q.headLock.Lock()
	node := q.head
	newHead := node.next
	if newHead == nil {
		q.headLock.Unlock()
		return
	}
	*value = newHead.value
	q.head = newHead
	q.headLock.Unlock()
}

type NodeAtomic struct {
	value int
	next  atomic.Pointer[NodeAtomic]
}

type QueueAtomic struct {
	head, tail atomic.Pointer[NodeAtomic]
}

func (q *QueueAtomic) Init() {
	node := &NodeAtomic{}
	q.head.Store(node)
	q.tail.Store(node)
}

func (q *QueueAtomic) Enqueue(value int) {
	node := &NodeAtomic{value: value}
	for {
		tail := q.tail.Load()
		next := tail.next.Load()
		if tail == q.tail.Load() {
			if next == nil {
				if tail.next.CompareAndSwap(nil, node) {
					q.tail.CompareAndSwap(tail, node)
					return
				}
			} else {
				q.tail.CompareAndSwap(tail, next)
			}
		}
	}
}

func (q *QueueAtomic) Dequeue(value *int) {
	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.next.Load()
		if head == q.head.Load() {
			if head == tail {
				if next == nil {
					return
				}
				q.tail.CompareAndSwap(tail, next)
			} else {
				*value = next.value
				if q.head.CompareAndSwap(head, next) {
					return
				}
			}
		}
	}
}

type QueueLock struct {
	head, tail         *Node
	headLock, tailLock sync.Mutex
}

func (q *QueueLock) Init() {
	node := &Node{}
	q.head, q.tail = node, node
}

func (q *QueueLock) Enqueue(value int) {
	node := &Node{value: value}
	q.tailLock.Lock()
	q.tail.next = node
	q.tail = node
	q.tailLock.Unlock()
}

func (q *QueueLock) Dequeue(value *int) bool {
	q.headLock.Lock()
	node := q.head
	newHead := node.next
	if newHead == nil {
		q.headLock.Unlock()
		return false
	}
	*value = newHead.value
	q.head = newHead
	q.headLock.Unlock()
	return true
}

func benchmarkQueue(q interface {
	Enqueue(int)
	Dequeue(*int)
}, numOps int) {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(i)
			var v int
			q.Dequeue(&v)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Time taken:", elapsed)
}

func main() {
	numOps := 100000
	fmt.Println("Benchmarking QueueMutex")
	qm := &QueueMutex{}
	qm.Init()
	benchmarkQueue(qm, numOps)

	fmt.Println("Benchmarking QueueAtomic")
	qa := &QueueAtomic{}
	qa.Init()
	benchmarkQueue(qa, numOps)
}
