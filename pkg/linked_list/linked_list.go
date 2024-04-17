package linked_list

import "sync"

type Element[T any] struct {
	data T
	next *Element[T]
}

func (e *Element[T]) Data() T {
	return e.data
}

type LinkedList[T any] struct {
	head    *Element[T]
	tail    *Element[T]
	current *Element[T]
	mu      *sync.RWMutex
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		mu: &sync.RWMutex{},
	}
}

func (ll *LinkedList[T]) Insert(data T) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	newEl := &Element[T]{data: data}

	if ll.head == nil {
		ll.head = newEl
		ll.tail = newEl
		ll.current = newEl
	} else {
		ll.tail.next = newEl
		ll.tail = newEl
	}
	ll.tail.next = ll.head
}

func (ll *LinkedList[T]) Next() *Element[T] {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.head == nil {
		return nil
	}
	current := ll.head
	ll.head = ll.head.next
	return current
}

func (ll *LinkedList[T]) GetElements() []T {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.head == nil {
		return nil
	}

	head := ll.head
	elements := []T{head.data}
	next := head.next

	for {
		if next == head {
			break
		}

		elements = append(elements, next.data)
		next = next.next
	}

	return elements
}
