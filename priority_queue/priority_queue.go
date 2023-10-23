package priorityqueue

import (
	"container/heap"
)

// Element - struct used to represent an item in the prioriry queue
type Element struct {
	Value    interface{}
	Priority int
	index    int
}

// maxheap - array of Element pointers that will be transformed in a heap
type maxheap []*Element

/** maxheap interface implementation **/

func (heap maxheap) Len() int {
	return len(heap)
}

func (heap maxheap) Less(i, j int) bool {
	return heap[i].Priority > heap[j].Priority
}

func (heap maxheap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
	heap[i].index = j
	heap[j].index = i
}

// Push - add an element at the end of the heap
func (heap *maxheap) Push(x interface{}) {
	n := len(*heap)
	item := x.(*Element)
	item.index = n
	*heap = append(*heap, item)
}

// Pop - pop the last element from the heap
func (heap *maxheap) Pop() interface{} {
	old := *heap
	n := len(old)

	item := old[n-1]
	old[n-1] = nil
	item.index = -1

	*heap = old[0 : n-1]

	return item
}

/** Wrapper over the go heap type **/

// PriorityQueue - wrapper over the heap implementation
type PriorityQueue struct {
	heap maxheap
}

// NewPriorityQueue - creates a new priority queue and initializes the heap
func NewPriorityQueue() PriorityQueue {
	priorityQueue := PriorityQueue{}
	priorityQueue.heap = make(maxheap, 0)
	heap.Init(&priorityQueue.heap)
	return priorityQueue
}

// Push - add a new element in the right position in the heap according to
// its priority
func (queue *PriorityQueue) Push(element *Element) {
	heap.Push(&queue.heap, element)
}

// Peek - returns the element with the highest priority from the heap
func (queue *PriorityQueue) Peek() *Element {
	return queue.At(0)
}

// Pop - removes and returns the element with the highest priority from the heap
func (queue *PriorityQueue) Pop() *Element {
	return heap.Pop(&queue.heap).(*Element)
}

// PopLowest - removes and returns the element with the lowest priority from the heap
func (queue *PriorityQueue) PopLowest() *Element {
	size := queue.Len()
	minIndex := size / 2

	for i := size/2 + 1; i < size; i++ {
		if queue.At(i).Priority < queue.At(minIndex).Priority {
			minIndex = i
		}
	}

	return queue.Remove(minIndex)
}

// Remove - removes and returns the element at position index from the heap
func (queue *PriorityQueue) Remove(index int) *Element {
	return heap.Remove(&queue.heap, index).(*Element)
}

// At - get the element at the specified index from the heap
func (queue *PriorityQueue) At(index int) *Element {
	return queue.heap[index]
}

// Len - return the number of elements in the heap
func (queue PriorityQueue) Len() int {
	return queue.heap.Len()
}

// IsEmpty - returns true if the heap has no elements, false otherwise
func (queue PriorityQueue) IsEmpty() bool {
	return queue.heap.Len() == 0
}