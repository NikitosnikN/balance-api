package nodepool

import "sync"

type ListElement struct {
	node *Node
	next *ListElement
}

func (e *ListElement) Node() *Node {
	return e.node
}

type LinkedList struct {
	head    *ListElement
	tail    *ListElement
	current *ListElement
	mu      *sync.RWMutex
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		mu: &sync.RWMutex{},
	}
}

func (ll *LinkedList) Insert(node *Node) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	newEl := &ListElement{node: node}

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

func (ll *LinkedList) Next() *ListElement {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.head == nil {
		return nil
	}
	current := ll.head
	ll.head = ll.head.next
	return current
}

func (ll *LinkedList) NextAliveNode() *Node {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.current == nil {
		return nil
	}

	current := ll.current
	next := current.next

	for {
		if next.node.IsAlive {
			ll.current = next
			break
		}

		if current == next {
			return nil
		}

		next = next.next
	}

	return ll.current.node
}

func (ll *LinkedList) GetElements() []*Node {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.head == nil {
		return nil
	}

	head := ll.head
	elements := []*Node{head.node}
	next := head.next

	for {
		if next == head {
			break
		}

		elements = append(elements, next.node)
		next = next.next
	}

	return elements
}
