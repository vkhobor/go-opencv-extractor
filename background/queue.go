package background

type BoundedQueue[T comparable] struct {
	items []T
	limit int
}

func NewBoundedQueue[T comparable](capacity int) BoundedQueue[T] {
	return BoundedQueue[T]{items: make([]T, 0, capacity), limit: capacity}
}

func (q *BoundedQueue[T]) Push(item T) {
	if len(q.items) == q.limit {
		q.items = q.items[1:]
	}
	q.items = append(q.items, item)

}

func (q *BoundedQueue[T]) Some(item T) bool {
	for _, i := range q.items {
		if i == item {
			return true
		}
	}

	return false
}

func (q *BoundedQueue[T]) SomeBy(predicate func(item T) bool) bool {
	for _, item := range q.items {
		if predicate(item) {
			return true
		}
	}

	return false
}
