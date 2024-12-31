package graph

import (
	"container/heap"
	"errors"
	"golang.org/x/exp/constraints"
)

// Priority queue implementation adapted from https://github.com/jupp0r/go-priority-queue

// PriorityQueue represents the queue
type PriorityQueue[T any, P constraints.Ordered] struct {
	itemHeap *itemHeap[T, P]
	lookup   map[interface{}]*item[T, P]
}

type itemHeap[T any, P constraints.Ordered] []*item[T, P]

type item[T any, P constraints.Ordered] struct {
	value    T
	priority P
	index    int
}

// checkInit initializes member variables, if needed
func (p *PriorityQueue[T, P]) checkInit() {
	if p.itemHeap == nil {
		p.itemHeap = &itemHeap[T, P]{}
	}
	if p.lookup == nil {
		p.lookup = make(map[interface{}]*item[T, P])
	}
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue[T, P]) Len() int {
	p.checkInit()
	return p.itemHeap.Len()
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *PriorityQueue[T, P]) Insert(v T, priority P) {
	p.checkInit()
	_, ok := p.lookup[v]
	if ok {
		return
	}
	newItem := &item[T, P]{
		value:    v,
		priority: priority,
	}
	heap.Push(p.itemHeap, newItem)
	p.lookup[v] = newItem
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *PriorityQueue[T, P]) Pop() (T, error) {
	p.checkInit()
	if len(*p.itemHeap) == 0 {
		zv := new(T)
		return *zv, errors.New("empty queue")
	}
	v := heap.Pop(p.itemHeap).(*item[T, P])
	delete(p.lookup, v.value)
	return v.value, nil
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *PriorityQueue[T, P]) UpdatePriority(x interface{}, newPriority P) {
	p.checkInit()
	v, ok := p.lookup[x]
	if !ok {
		return
	}
	v.priority = newPriority
	heap.Fix(p.itemHeap, v.index)
}

func (ih *itemHeap[T, P]) Len() int {
	return len(*ih)
}

func (ih *itemHeap[T, P]) Less(i, j int) bool {
	return (*ih)[i].priority < (*ih)[j].priority
}

func (ih *itemHeap[T, P]) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeap[T, P]) Push(x interface{}) {
	it := x.(*item[T, P])
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeap[T, P]) Pop() interface{} {
	old := *ih
	item := old[len(old)-1]
	*ih = old[0 : len(old)-1]
	return item
}
