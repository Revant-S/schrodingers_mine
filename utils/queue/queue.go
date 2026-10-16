package queue

import "fmt"

type Queue[T any] struct {
	container []T
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		container: make([]T, 0),
	}
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(item T) {
	q.container = append(q.container, item)
}

// Pop removes and returns the element from the front of the queue.
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if len(q.container) == 0 {
		return zero, false
	}
	item := q.container[0]
	q.container[0] = zero
	q.container = q.container[1:]
	return item, true
}

func (q *Queue[T]) Front() (T, bool) {
	if len(q.container) == 0 {
		var zero T
		return zero, false
	}
	return q.container[0], true
}

// IsEmpty checks if the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.container) == 0
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	return len(q.container)
}

// Print displays the elements of the queue from front to back.
func (q *Queue[T]) Print() {
	fmt.Println(q.container)
}
