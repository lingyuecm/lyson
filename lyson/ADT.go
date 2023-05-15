package lyson

type BiLinkedNode[T any] struct {
	data         T
	previousNode *BiLinkedNode[T]
	nextNode     *BiLinkedNode[T]
}

type Stack[T any] struct {
	peek *BiLinkedNode[T]
	size int
}

func NewStack[T any]() *Stack[T] {
	s := new(Stack[T])

	s.peek = new(BiLinkedNode[T])
	s.size = 0

	return s
}

func (s *Stack[T]) Push(element T) {
	node := new(BiLinkedNode[T])
	node.data = element

	node.previousNode = s.peek
	s.peek.nextNode = node

	s.peek = node
	s.size = s.size + 1
}

func (s *Stack[T]) Pop() T {
	if 0 == s.Size() {
		panic("Empty Stack")
	}
	node := s.peek
	s.peek = node.previousNode

	node.previousNode = nil
	s.peek.nextNode = nil
	s.size = s.size - 1

	return node.data
}

func (s *Stack[T]) Peek() T {
	if 0 == s.size {
		panic("Empty Stack")
	}
	node := s.peek
	return node.data
}

func (s *Stack[T]) Size() int {
	return s.size
}
