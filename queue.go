package valize

type QueueStrategy interface {
	Push(elem []byte) error
	Peek() ([]byte, error)
	Pop() ([]byte, error)
	Clear() error
}

type Queue struct {
	Backend QueueStrategy
}

func (s *Queue) Push(elem []byte) error {
	s.Backend.Push(elem)
	return nil
}
func (s *Queue) Peek() ([]byte, error) {
	return s.Backend.Peek()
}
func (s *Queue) Pop() ([]byte, error) {
	return s.Backend.Pop()
}
func (s *Queue) Cleear() error {
	return s.Backend.Clear()
}
