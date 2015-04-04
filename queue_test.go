package valize

import (
	"errors"
	"testing"
)

type MockQueue struct {
	Store [][]byte
}

func NewMockQueue() *MockQueue {
	return &MockQueue{Store: [][]byte{}}
}

func (s *MockQueue) Push(elem []byte) error {
	s.Store = append(s.Store, elem)
	return nil
}
func (s *MockQueue) Peek() ([]byte, error) {
	if len(s.Store) > 0 {
		return s.Store[0], nil
	}
	return nil, errors.New("Empty queue")
}
func (s *MockQueue) Pop() ([]byte, error) {
	if len(s.Store) > 0 {
		elem := s.Store[0]
		s.Store = s.Store[1:]
		return elem, nil
	}
	return nil, errors.New("Empty queue")
}
func (s *MockQueue) Clear() error {
	s.Store = nil
	return nil
}

func TestQueueInit(t *testing.T) {
	actual := &Queue{Backend: NewMockQueue()}
	if actual == nil {
		t.Error("MockQueue failed to initialize")
	}
}

func TestQueuePush(t *testing.T) {
	queue := &Queue{Backend: NewMockQueue()}
	if err := queue.Push([]byte("a")); err != nil {
		t.Error("Error Pushing to queue")
	}
	if err := queue.Push([]byte("b")); err != nil {
		t.Error("Error Pushing to queue")
	}
}

func TestQueuePeek(t *testing.T) {
	queue := &Queue{Backend: NewMockQueue()}

	if err := queue.Push([]byte("a")); err != nil {
		t.Error("Error Pushing to queue")
	}

	if err := queue.Push([]byte("b")); err != nil {
		t.Error("Error Pushing to queue")
	}

	actual1, err1 := queue.Peek()
	actual2, err2 := queue.Peek()
	actual3, err3 := queue.Peek()

	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Error Peeking from queue")
	}

	expected := []byte("a")
	if actual1[0] != expected[0] {
		t.Error("Assert Fail", actual1, expected)
	}
	if actual2[0] != expected[0] {
		t.Error("Assert Fail", actual2, expected)
	}
	if actual3[0] != expected[0] {
		t.Error("Assert Fail", actual3, expected)
	}
}

func TestQueuePop(t *testing.T) {
	queue := &Queue{Backend: NewMockQueue()}

	if err := queue.Push([]byte("a")); err != nil {
		t.Error("Error Pushing to queue")
	}

	if err := queue.Push([]byte("b")); err != nil {
		t.Error("Error Pushing to queue")
	}

	actual1, err1 := queue.Pop()
	actual2, err2 := queue.Pop()
	actual3, err3 := queue.Pop()

	expected1 := []byte("a")
	if actual1[0] != expected1[0] {
		t.Error("Assert Fail", actual1, expected1)
	}

	expected2 := []byte("b")
	if actual2[0] != expected2[0] {
		t.Error("Assert Fail", actual2, expected2)
	}

	if err1 != nil || err2 != nil || actual3 != nil {
		t.Error("Error Peeking from queue")
	}

	if err3 == nil {
		t.Error("Error queue is non-empty")
	}

}
